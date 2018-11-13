package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/bcext/gcash/wire"
	"github.com/qshuai/tcolor"
)

func main() {
	tx := flag.String("tx", "", "Please input a transaction with hexadecimal decoded")
	in := flag.Bool("input", false, "Whether fetch the specified transaction input or not(default false)")
	out := flag.Bool("output", false, "Whether fetch the specified transaction output or not(default false)")
	num := flag.Int("num", 0, "The specified item with offset(based 0)")
	flag.Parse()

	if *tx == "" {
		fmt.Println(tcolor.WithColor(tcolor.Red, "empty transaction string not allowed"))
		os.Exit(1)
	}

	b, err := hex.DecodeString(*tx)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "decode hexadecimal string failed: "+err.Error()))
		os.Exit(1)
	}

	var t wire.MsgTx
	err = t.Deserialize(bytes.NewReader(b))
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "decode transaction failed: "+err.Error()))
		os.Exit(1)
	}

	if *in {
		if *num > len(t.TxIn)-1 {
			fmt.Println(tcolor.WithColor(tcolor.Red, "the offset overflow"))
			os.Exit(1)
		}

		fmt.Println("previous outpoint:", t.TxIn[*num].PreviousOutPoint.String())
		fmt.Println("signature:", hex.EncodeToString(t.TxIn[*num].SignatureScript))
		fmt.Printf("sequence: %#x\n", t.TxIn[*num].Sequence)

		os.Exit(0)
	}

	if *out {
		if *num > len(t.TxOut)-1 {
			fmt.Println(tcolor.WithColor(tcolor.Red, "the offset overflow"))
			os.Exit(1)
		}

		fmt.Println("value:", t.TxOut[*num].Value)
		fmt.Println("PkScript:", hex.EncodeToString(t.TxOut[*num].PkScript))
	}
}
