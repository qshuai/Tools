package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	args := os.Args
	var length uint32
	if len(args) < 2 {
		length = rand.Uint32() % 64

		if length < 4 {
			length += 4
		}
	} else {
		ret, err := strconv.Atoi(args[1])
		if err != nil {
			panic("Please check the input number!")
		}

		if ret <= 0 || ret > 520 {
			panic("Input number should larger than 0 and less than 520")
		}

		length = uint32(ret)
	}

	data := make([]byte, length/2)
	rand.Read(data)
	fmt.Println(strings.ToUpper(hex.EncodeToString(data)))
}
