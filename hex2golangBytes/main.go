package main

import (
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/qshuai/tcolor"
)

func main() {
	data := flag.String("data", "", "please a hexadecimal string")
	flag.Parse()

	if len(*data) == 0 {
		fmt.Printf(tcolor.WithColor(tcolor.Red, "empty data not allowed"))
		return
	}

	b, err := hex.DecodeString(*data)
	if err != nil {
		panic(err)
	}

	for _, item := range b {
		if item <= 0x0f {
			fmt.Printf("0x0%x, ", item)
		} else {
			fmt.Printf("%#x, ", item)
		}
	}
}
