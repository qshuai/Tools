package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/bcext/gcash/wire"
	"github.com/qshuai/tcolor"
	"github.com/tidwall/gjson"
)

const (
	api = "https://developer-bch-chain.api.btc.com/appkey-2f7c183e3e9e/"
)

func main() {
	height := flag.Int("height", 100, "please input the specified block height")
	limit := flag.Int("limit", 100, "How many blocks do you want")
	flag.Parse()

	headers := make([]wire.BlockHeader, *limit)
	start := *height
	for i := *limit; i > 0; i-- {
		res, err := http.Get(api + strconv.Itoa(start))
		if err != nil || res.StatusCode != http.StatusOK {
			fmt.Println(tcolor.WithColor(tcolor.Red, "http request failed: "+err.Error()))
			os.Exit(1)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "read http response failed: "+err.Error()))
			os.Exit(1)
		}

		content := string(b)
		data := gjson.Get(content, "data")

		var header wire.BlockHeader
		header.Timestamp = time.Unix(data.Get("timestamp").Int(), 0)
		prevHash, err := chainhash.NewHashFromStr(data.Get("prev_block_hash").String())
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "decode block hash from string failed: "+err.Error()))
			os.Exit(1)
		}
		header.PrevBlock = *prevHash
		header.Version = int32(data.Get("version").Int())
		header.Bits = uint32(data.Get("bits").Int())
		header.Nonce = uint32(data.Get("nonce").Int())

		merkleHash, err := chainhash.NewHashFromStr(data.Get("mrkl_root").String())
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "decode block hash from string failed: "+err.Error()))
			os.Exit(1)
		}
		header.MerkleRoot = *merkleHash
		headers[i-1] = header

		start--

		// prevent from abuse the API server
		time.Sleep(200 * time.Millisecond)
	}

	formatHeaders(*height-100+1, headers)
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
