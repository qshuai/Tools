package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/qshuai/tcolor"
	"os"
)

func main() {
	pubkey := flag.String("pubkey", "", "a compressed or uncompressed public key")
	flag.Parse()

	if pubkey == nil || len(*pubkey) == 0 {
		fmt.Println("please a public key string(hexadecimal)")
		os.Exit(1)
	}

	pkBytes, err := hex.DecodeString(*pubkey)
	if err != nil {
		panic(err)
	}

	pk, err := btcec.ParsePubKey(pkBytes, btcec.S256())
	if err != nil {
		panic(err)
	}

	pkHash := btcutil.Hash160(pk.SerializeCompressed())

	addr, err := btcutil.NewAddressPubKeyHash(pkHash, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}

	fmt.Println(tcolor.WithColor(tcolor.Green, addr.String()))
}
