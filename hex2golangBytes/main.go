package main

import (
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/qshuai/tcolor"
)

func main() {
	data := flag.String("data", "", "please a hexadecimal string")
	num := flag.Int("num", 12, "bytes per line")
	flag.Parse()

	if len(*data) == 0 {
		fmt.Printf(tcolor.WithColor(tcolor.Red, "empty data not allowed"))
		return
	}

	b, err := hex.DecodeString(*data)
	if err != nil {
		panic(err)
	}

	var lineOffset int
	fmt.Print("[]byte{")
	for index, item := range b {
		if item <= 0x0f {
			fmt.Printf("0x0%x", item)
		} else {
			fmt.Printf("%#x", item)
		}

		if index != len(b)-1 {
			fmt.Print(", ")
		}

		lineOffset++
		if lineOffset%*num == 0 {
			fmt.Print("\n")
		}

		if (index < *num-2) && lineOffset%(*num-2) == 0 {
			fmt.Print("\n")
			lineOffset = 0
		}
	}
	fmt.Print("}")
}
