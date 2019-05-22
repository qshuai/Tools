package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/wire"
	"github.com/qshuai/tcolor"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println(tcolor.WithColor(tcolor.Red, "invalid parameter[double256hash msg]"))
		os.Exit(1)
	}

	msg, err := hex.DecodeString(args[1])
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "invalid hexadecimal string"))
		os.Exit(1)
	}

	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(msg))
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "invalid transaction raw hexadecimal string"))
		os.Exit(1)
	}

	fmt.Println(tcolor.WithColor(tcolor.Green, "txid:"+tx.TxHash().String()))
	fmt.Println(tcolor.WithColor(tcolor.Green, "hash:"+tx.WitnessHash().String()))
	//
	//first := sha256.Sum256(msg)
	//second := sha256.Sum256(first[:])
	//hash, err := chainhash.NewHash(second[:])
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(tcolor.WithColor(tcolor.Green, hash.String()))
	//
	//tx := wire.MsgTx{}
	//tx.TxHash()

}
