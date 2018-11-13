package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bcext/gcash/rpcclient"
	"github.com/bcext/gcash/wire"
	"github.com/qshuai/tcolor"
)

var (
	client *rpcclient.Client
)

func main() {
	height := flag.Int("height", 100, "Please input the specified block height")
	host := flag.String("rpchost", "127.0.0.1:8332", "Please input rpc host(ip:port)")
	user := flag.String("rpcuser", "", "Please input your rpc username")
	passwd := flag.String("rpcpassword", "", "Please input your rpc password")
	limit := flag.Int("limit", 100, "How many blocks do you want")
	flag.Parse()

	client = GetRPC(*host, *user, *passwd)

	headers := make([]wire.BlockHeader, *limit)
	start := *height
	for i := *limit; i > 0; i-- {
		hash, err := client.GetBlockHash(int64(start))
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "get block hash failed: "+err.Error()))
			os.Exit(1)
		}

		block, err := client.GetBlock(hash)
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "get block failed: "+err.Error()))
			os.Exit(1)
		}

		headers[i-1] = block.Header

		start--
	}

	formatHeaders(*height-100+1, headers)
}

func GetRPC(host, user, passwd string) *rpcclient.Client {
	if client != nil {
		return client
	}

	// rpc client instance
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         passwd,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	c, err := rpcclient.New(connCfg, nil)
	if err != nil {
		panic(err)
	}

	client = c
	return c
}

func formatHeaders(start int, headers []wire.BlockHeader) {
	for _, header := range headers {
		fmt.Println("> block header:", start)
		fmt.Println("\tversion:      ", header.Version)
		fmt.Println("\tprevBlockHash:", header.PrevBlock)
		fmt.Println("\tMerkleRoot:   ", header.MerkleRoot)
		fmt.Println("\tBits:         ", header.Bits)
		fmt.Println("\tNonce:        ", header.Nonce)
		fmt.Println("\tTimestamp:    ", header.Timestamp.Unix())

		start++

		// blank line
		fmt.Println()
	}
}
