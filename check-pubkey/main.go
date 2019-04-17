package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/bcext/gcash/btcec"
	"github.com/qshuai/tcolor"
)

func main() {
	pk := flag.String("pk", "", "input a compressed or uncompressed public key")
	flag.Parse()

	if pk == nil || *pk == "" {
		fmt.Println("please a input public key string")
		os.Exit(1)
	}

	pubkey, err := hex.DecodeString(*pk)
	if err != nil {
		fmt.Println("invalid hexadecimal string for public key")
		os.Exit(1)
	}

	_, err = btcec.ParsePubKey(pubkey, btcec.S256())
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "invalid"))
	} else {
		fmt.Println(tcolor.WithColor(tcolor.Green, "valid"))
	}
}
