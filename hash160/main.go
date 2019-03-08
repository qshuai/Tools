package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/qshuai/tcolor"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println(tcolor.WithColor(tcolor.Red, "check argument"))
		os.Exit(1)
	}

	data, err := hex.DecodeString(args[1])
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "invalid hexadecimal string"))
		os.Exit(1)
	}

	fmt.Println(hex.EncodeToString(btcutil.Hash160(data)))
}
