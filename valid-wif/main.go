package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg"
	"github.com/qshuai/tcolor"
)

var net = map[string]*chaincfg.Params{
	"mainnet": &chaincfg.MainNetParams,
	"testnet": &chaincfg.TestNet3Params,
	"regtest": &chaincfg.RegressionNetParams,
}

func main() {
	privkey := flag.String("privkey", "", "please input a wif private key")
	chain := flag.String("param", "mainnet", "select one item from mainnet/testnet/regtest")
	flag.Parse()

	param, ok := net[*chain]
	if !ok {
		fmt.Println(tcolor.WithColor(tcolor.Red, *chain+" not existed, should select from mainnet/testnet/regtest"))
		return
	}

	wif, err := cashutil.DecodeWIF(*privkey)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "invalid private key: "+err.Error()))
		os.Exit(1)
	}

	// get bitcoin address for sender
	pubKey := wif.PrivKey.PubKey()
	pubKeyBytes := pubKey.SerializeCompressed()
	pkHash := cashutil.Hash160(pubKeyBytes)
	addr, err := cashutil.NewAddressPubKeyHash(pkHash, param)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "address encode failed, please check your privkey: "+err.Error()))
		os.Exit(1)
	}

	fmt.Println("address:")
	fmt.Println("  base58 encoded address:", tcolor.WithColor(tcolor.Green, addr.EncodeAddress(false)))
	fmt.Println("  bech32 encoded address:", tcolor.WithColor(tcolor.Green, addr.EncodeAddress(true)))
}
