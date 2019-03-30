package main

import (
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/qshuai/tcolor"
)

func main() {
	rawtx := flag.String("rawtx", "", "Please input the raw hexadecimal transaction string")
	flag.Parse()
	if rawtx == nil || len(*rawtx) == 0 {
		fmt.Println(tcolor.WithColor(tcolor.Red, "lack of raw hexadecimal transaction string"))
		return
	}

	b, err := hex.DecodeString(*rawtx)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "not valid hexadecimal transaction string"))
		return
	}

	txHash := chainhash.DoubleHashH(b)
	fmt.Println(tcolor.WithColor(tcolor.Green, "txhash: "+txHash.String()))
}
