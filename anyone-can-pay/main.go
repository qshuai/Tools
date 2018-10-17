package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg"
	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/bcext/gcash/txscript"
	"github.com/bcext/gcash/wire"
	"github.com/qshuai/tcolor"
)

func main() {
	hash := flag.String("hash", "", "Please input the anyonecanpay transaction hash")
	vout := flag.Int("vout", 0, "Please input the anyonecanpay output index")
	redeem := flag.String("redeem", "", "Please input a redeem address for you")
	amount := flag.Int("amount", 0, "Please input redeem amount")
	flag.Parse()

	var tx wire.MsgTx
	prevHash, _ := chainhash.NewHashFromStr(*hash)
	outpoint := wire.NewOutPoint(prevHash, uint32(*vout))
	txin := wire.NewTxIn(outpoint, nil)
	txin.Sequence = 0xffffffff
	tx.TxIn = append(tx.TxIn, txin)

	address, err := cashutil.DecodeAddress(*redeem, &chaincfg.TestNet3Params)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Please check your address validation"))
	}

	pkscript, _ := txscript.PayToAddrScript(address)
	out := wire.NewTxOut(int64(*amount), pkscript)
	tx.TxOut = append(tx.TxOut, out)
	tx.LockTime = 0

	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	tx.Serialize(buf)

	fmt.Println(hex.EncodeToString(buf.Bytes()))
}
