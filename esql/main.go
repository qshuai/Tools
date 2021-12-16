package main

import (
	"flag"
	"fmt"
	"os"
	"unsafe"

	"github.com/cch123/elasticsql"
	"github.com/tidwall/pretty"
)

func main() {
	sql := flag.String("sql", "", "sql statement")
	flag.Parse()

	if len(*sql) <= 0 {
		fmt.Println("please input sql statement!!")
		os.Exit(1)
	}

	dsl, _, err := elasticsql.Convert(*sql)
	if err != nil {
		panic(err)
	}

	bs := *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{dsl, len(dsl)}))

	bs = pretty.Pretty(bs)
	bs = pretty.Color(bs, pretty.TerminalStyle)
	fmt.Println(*(*string)(unsafe.Pointer(&bs)))
}
