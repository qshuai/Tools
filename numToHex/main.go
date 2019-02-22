package main

import (
	"fmt"
	"github.com/qshuai/tcolor"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println(tcolor.WithColor(tcolor.Red, "please input one number at least"))
		os.Exit(1)
	}

	for _, num := range args[1:] {
		if strings.Contains(num, ".") {
			fmt.Println(tcolor.WithColor(tcolor.Red, "do not support decimal"))
			os.Exit(1)
		}

		n, err := strconv.ParseUint(num, 10, 64)
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "convert the argument to number failed: " + err.Error()))
			os.Exit(1)
		}
		fmt.Printf("0x%x\n", n)
	}
}
