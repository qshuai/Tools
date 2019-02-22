package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/qshuai/tcolor"
	"os"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println(tcolor.WithColor(tcolor.Red, "please input a hexadecimal string"))
		os.Exit(1)
	}

	for _, str := range args[1:] {
		bs, err := hex.DecodeString(str)
		if err != nil {
			panic(err)
		}
		if len(bs) > 8 {
			fmt.Println(tcolor.WithColor(tcolor.Red, "integer too big"))
			os.Exit(1)
		}


		fmt.Println(">")
		{
			r := make([]byte, 8)
			copy(r[8-len(bs):], bs)
			bigNum := binary.BigEndian.Uint64(r)
			fmt.Println("  BigEndian:", bigNum)
		}

		{
			r := make([]byte, 8)
			copy(r, bs)
			littleNum := binary.LittleEndian.Uint64(r)
			fmt.Println("  LittleEndian:", littleNum)
		}

		fmt.Println()
	}
}
