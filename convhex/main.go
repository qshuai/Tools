package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	input := flag.String("hex", "", "please input hexadecimal string")
	flag.Parse()
	if input == nil {
		panic("not input")
	}

	if len(*input)%2 != 0 {
		panic("incorrect hexadecimal string")
	}

	origin := *input

	var ret strings.Builder
	for i := len(origin) - 1; i > 0; {
		_, err := ret.WriteString(origin[i-1 : i+1])
		if err != nil {
			panic(err)
		}
		i = i - 2
	}

	fmt.Println(ret.String())
}
