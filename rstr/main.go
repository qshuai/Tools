package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"strconv"

	"github.com/bcext/cashutil/base58"
	"github.com/qshuai/tcolor"
)

type encodeType string

const (
	base58Encoding encodeType = "base58"
	base64Encoding encodeType = "base64"

	bytesLengthLimit   = 500
	defaultBytesLength = 60
)

func main() {
	length := flag.Int("length", defaultBytesLength, "Please input the needed length of the generated random string (base58 encoded string length is not precision)")
	encode := flag.String("encode", string(base64Encoding), "Please input encode type (base64 encoded string with character '+/=')")
	flag.Parse()

	if *length > bytesLengthLimit {
		fmt.Println(tcolor.WithColor(tcolor.Yellow, "the input length too big. use default length: "+strconv.Itoa(defaultBytesLength)))
		*length = defaultBytesLength
	}

	data := make([]byte, int(*length*3)/4)
	rand.Read(data)
	if *encode == string(base58Encoding) {
		fmt.Println(tcolor.WithColor(tcolor.Green, base58.Encode(data)))
		return
	}

	fmt.Println(tcolor.WithColor(tcolor.Green, base64.StdEncoding.EncodeToString(data)))
}
