package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"decimal"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg"
	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/bcext/gcash/txscript"
	"github.com/bcext/gcash/wire"
	"github.com/qshuai/tcolor"
	"github.com/tidwall/gjson"
)

const (
	bitcoinCashAPI  = "https://bch-tchain.api.btc.com/v3"
	defaultPageSize = 50

	// address pair
	bech32Address = "bchtest:qqwpvaha3leydercn7kckkuh9ufxaplcmsn48e8v3m"
	base58Address = "mi5U8JnLMLiVrms3mW9YNvz5nSYC57Q7G9"

	privkey = "cRL6HJZYSF1JMUSyP6PsKMRD9PhCS1acUSoKWh9p5Bf5iY4SPq5j"

	// One transaction(total P2PKH) with one input and one output,
	// its size is 192 bytes.
	// 0.1 BCH utxo will create max feerate: 0.52080489 BCH/KB
	minUtxoValue = 10000000

	defaultSignatureSize = 107
	defaultSequence      = 0xffffffff
)

func main() {
	feerate := flag.Float64("feerate", 0.00001, "Please input the specified feerate to create transaction")
	flag.Parse()

	utxos, err := getUnspent(base58Address, 1)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Request utxo list failed: "+err.Error()))
		return
	}

	list := gjson.Get(utxos, "data.list").Array()
	var targetIdx int
	for index, item := range list {
		if item.Get("value").Int() >= minUtxoValue {
			targetIdx = index
			break
		}
	}

	// decode address
	address, _ := cashutil.DecodeAddress(bech32Address, &chaincfg.TestNet3Params)

	// parse privkey
	wif, err := cashutil.DecodeWIF(privkey)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Privkey format error"))
		os.Exit(1)
	}

	inputValue, tx, err := assembleTx(list[targetIdx], address, wif, *feerate)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Assemble transaction or sign error:"+err.Error()))
		os.Exit(1)
	}

	var positive, negative bool
again:
	realFeeRate := calcFeeRate(tx.SerializeSize(), inputValue, tx.TxOut[0].Value)
	if realFeeRate < *feerate {
		tx.TxOut[0].Value--
		if positive {
			goto next
		}
		negative = true
		goto again
	} else {
		tx.TxOut[0].Value++
		if negative {
			goto next
		}
		positive = true
		goto again
	}
next:

	buf := bytes.NewBuffer(nil)
	err = tx.Serialize(buf)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Transaction serialize error:"+err.Error()))
		os.Exit(1)
	}

	// output result
	fmt.Println("txhash:         ", tcolor.WithColor(tcolor.Green, tx.TxHash().String()))
	fmt.Println("raw transaction:", tcolor.WithColor(tcolor.Green, hex.EncodeToString(buf.Bytes())))
}

// get raw string of unspent list for the specified address
func getUnspent(addr string, page int) (string, error) {
	url := bitcoinCashAPI + "/address/" + addr + "/unspent?pagesize=" +
		strconv.Itoa(defaultPageSize) + "&page=" + strconv.Itoa(page)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", errors.New("request failed")
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func assembleTx(utxo gjson.Result, address cashutil.Address, wif *cashutil.WIF, feerate float64) (int64, *wire.MsgTx, error) {
	var tx wire.MsgTx
	tx.Version = 1
	tx.LockTime = 0

	tx.TxOut = make([]*wire.TxOut, 1)
	pkScript, _ := txscript.PayToAddrScript(address)
	tx.TxOut[0] = &wire.TxOut{PkScript: pkScript}

	hash, err := chainhash.NewHashFromStr(utxo.Get("tx_hash").String())
	if err != nil {
		return 0, nil, err
	}

	outpoint := wire.NewOutPoint(hash, uint32(utxo.Get("tx_output_n").Int()))
	tx.TxIn = append(tx.TxIn, wire.NewTxIn(outpoint, nil))
	tx.TxIn[0].Sequence = defaultSequence

	inputTotal := utxo.Get("value").Int()
	txsize := tx.SerializeSize() + defaultSignatureSize

	fee := decimal.NewFromFloat(feerate * 1e5).Mul(decimal.New(int64(txsize), 0)).Truncate(0).IntPart()

	outValue := inputTotal - fee
	tx.TxOut[0].Value = outValue

	// sign the transaction
	rawtx, err := sign(&tx, []int64{utxo.Get("value").Int()}, pkScript, wif)
	return inputTotal, rawtx, err
}

func sign(tx *wire.MsgTx, inputValueSlice []int64, pkScript []byte, wif *cashutil.WIF) (*wire.MsgTx, error) {
	for idx, _ := range tx.TxIn {
		sig, err := txscript.RawTxInSignature(tx, idx, pkScript, cashutil.Amount(inputValueSlice[idx]),
			txscript.SigHashAll|txscript.SigHashForkID, wif.PrivKey)
		if err != nil {
			return nil, err
		}
		sig, err = txscript.NewScriptBuilder().AddData(sig).Script()
		if err != nil {
			return nil, err
		}
		pk, err := txscript.NewScriptBuilder().AddData(wif.PrivKey.PubKey().SerializeCompressed()).Script()
		if err != nil {
			return nil, err
		}
		sig = append(sig, pk...)
		tx.TxIn[0].SignatureScript = sig

		engine, err := txscript.NewEngine(pkScript, tx, idx, txscript.StandardVerifyFlags,
			nil, nil, inputValueSlice[idx])
		if err != nil {
			return nil, err
		}

		// verify the signature
		err = engine.Execute()
		if err != nil {
			return nil, err
		}
	}

	return tx, nil
}

func calcFeeRate(txsize int, inputValue, outputValue int64) float64 {
	fee := inputValue - outputValue
	feeRate, _ := decimal.New(fee, 0).Div(decimal.New(int64(txsize), 0)).Mul(decimal.New(1, -5)).Truncate(8).Float64()
	return feeRate
}
