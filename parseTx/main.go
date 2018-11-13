package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/bcext/gcash/wire"
	"github.com/davecgh/go-spew/spew"
	"github.com/qshuai/tcolor"
)

func main() {
	tx := flag.String("tx", "", "Please input a transaction with hexadecimal decoded")
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

	spew.Dump(t)
}
