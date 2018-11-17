package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg"
	"github.com/bcext/gcash/rpcclient"
	"github.com/qshuai/tcolor"
	"github.com/shopspring/decimal"
)

const (
	maxRetries = 10
)

var params = map[string]*chaincfg.Params{
	"mainnet": &chaincfg.MainNetParams,
	"testnet":&chaincfg.TestNet3Params,
	"regtest": &chaincfg.RegressionNetParams,
}

func main() {
	chain := flag.String("chainparam", "testnet", "blockchain parameter")
	host := flag.String("rpchost", "127.0.0.1:18332", "host for rpc server[ip:port]")
	user := flag.String("rpcuser", "", "rpc username for rpc server auth")
	passwd := flag.String("rpcpasswd", "", "rpc password for rpc server auth")

	address := flag.String("address", "", "bitcoin address for receiver")
	amount := flag.Float64("amount", 0.1, "amount sending for each transaction[in bitcoin]")
	count := flag.Int("count", 50, "create how many transactions")
	flag.Parse()

	if *user == "" || *passwd == "" || *address == "" {
		fmt.Println(tcolor.WithColor(tcolor.Red, "empty user, password and address not allowd"))
		os.Exit(1)
	}

	param := getParam(*chain)
	if param == nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "chainparam not existed"))
		os.Exit(1)
	}

	conf := rpcclient.ConnConfig{
		Host:*host,
		User:*user,
		Pass:*passwd,
		HTTPPostMode:true,
		DisableTLS:true,
	}
	client,err := rpcclient.New(&conf, nil)
	if err != nil {
		fmt.Print(tcolor.WithColor(tcolor.Red, "rpc client instance initial failed"))
		os.Exit(1)
	}

	// decode address
	addr, err := cashutil.DecodeAddress(*address, param)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "invalid bitcoin address"))
		os.Exit(1)
	}

	// convert amount
	value := decimal.NewFromFloat(*amount).Mul(decimal.New(1e8,0)).IntPart()
	var failed int
	for i := 0; i < *count; i++ {
		hash, err := client.SendToAddress(addr,cashutil.Amount(value))
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "send failed"))

			failed++
			if failed >= maxRetries {
				fmt.Println(tcolor.WithColor(tcolor.Red, "exit: failed to many times"))
				os.Exit(1)
			}
		} else {
			fmt.Println(tcolor.WithColor(tcolor.Green, hash.String()))
		}
	}

	fmt.Println(tcolor.WithColor(tcolor.Cyan, "Done!"))
}

func getParam(param string) *chaincfg.Params {
	chain , ok := params[param]
	if !ok {
		return nil
	}

	return chain
}



