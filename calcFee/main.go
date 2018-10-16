package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"decimal"

	"github.com/bcext/gcash/wire"
	"github.com/pkg/errors"
	"github.com/qshuai/tcolor"
	"github.com/tidwall/gjson"
)

const (
	bitcoinCashAPI = "https://bch-tchain.api.btc.com/v3"
)

func main() {
	rawtx := flag.String("rawtx", "", "Please input the raw hexadecimal transaction string")
	flag.Parse()
	if rawtx == nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "lack of raw hexadecimal transaction string"))
		return
	}

	txBytes, err := hex.DecodeString(*rawtx)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "invalid hexadecimal string"))
		return
	}

	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewBuffer(txBytes))
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "transaction deserialize failed"))
		return
	}

	tx.TxHash()

	// calculate fee
	var inputValue, outputValue int64
	var unconfirmedArray []string
	for _, in := range tx.TxIn {
		value, unconfirmed, err := getInputValue(in.PreviousOutPoint.Hash.String(), in.PreviousOutPoint.Index)
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "calculate input total value err: "+err.Error()))
			return
		}

		if unconfirmed {
			unconfirmedArray = append(unconfirmedArray, in.PreviousOutPoint.Hash.String()+":"+
				strconv.Itoa(int(in.PreviousOutPoint.Index)))

			continue
		}

		inputValue += value
	}

	for _, out := range tx.TxOut {
		outputValue += out.Value
	}

	// So we can get fee and feerate, or unconfirmed utxo if encounter
	if unconfirmedArray != nil {
		fmt.Println(tcolor.WithColor(tcolor.Yellow, "unconfirmed utxo in this transaction as following:"))
		for _, item := range unconfirmedArray {
			fmt.Println(tcolor.WithColor(tcolor.Green, item))
		}

		return
	}

	fee := inputValue - outputValue
	txSize := len(*rawtx) / 2
	feeRateWithSatoshi := decimal.New(fee, 0).Div(decimal.New(int64(txSize), 0)).IntPart()
	feeRateWithBCH, _ := decimal.New(fee, 0).Div(decimal.New(int64(txSize), 0)).Mul(decimal.New(1, -5)).Truncate(8).Float64()

	fmt.Println(tcolor.WithColor(tcolor.Green, fmt.Sprintf("fee    : %s", decimal.New(fee, 0).Mul(decimal.New(1, -8)).Truncate(8).String())))
	fmt.Println(tcolor.WithColor(tcolor.Green, fmt.Sprintf("feeRate: %d Satoshi/Byte", feeRateWithSatoshi)))
	fmt.Println(tcolor.WithColor(tcolor.Green, fmt.Sprintf("feeRate: %f BCH/KB", feeRateWithBCH)))
}

func getInputValue(hash string, vout uint32) (int64, bool, error) {
	res, err := http.Get(bitcoinCashAPI + "/tx/" + hash)
	if err != nil {
		return 0, false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, false, errors.New("request failed with http failed code: " + strconv.Itoa(res.StatusCode))
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, false, err
	}

	// check whether this utxo is uncofirmed or not
	if gjson.Get(string(b), "err_no").Int() == 1 {
		return 0, true, nil
	}

	outs := gjson.Get(string(b), "data.outputs").Array()
	if int(vout) >= len(outs) {
		return 0, false, errors.New("output index is overflow")
	}

	return outs[vout].Get("value").Int(), false, nil
}
