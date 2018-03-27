//this program will calculate the length of input string.
//you can build a executable file and move it to your path.
//usage: go build -o len main.go
//commandline: len input your string

package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	length := len(args)
	if length < 2 {
		fmt.Println("   \033[0;31mplease input a or more string\n\033[0m")
		fmt.Println("   usage:\033[0;32mlen hello world test usage\033[0m")
		fmt.Println("   output:\033[0;32m")
		fmt.Println("   \033[0;32m       5")
		fmt.Println("   \033[0;32m       5")
		fmt.Println("   \033[0;32m       4")
		fmt.Println("   \033[0;32m       5")
		return
	}
	for i := 1; i < length; i++ {
		fmt.Print("\033[0;32m ")
		fmt.Println(len(args[i]))
	}
}