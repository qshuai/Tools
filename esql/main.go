package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cch123/elasticsql"
)

func main() {
	sql := flag.String("sql", "", "sql statement")
	flag.Parse()

	if len(*sql) <= 0 {
		fmt.Println("please input sql statement!!")
		os.Exit(1)
	}

	dsl, esType, err := elasticsql.Convert(sql)
	if err != nil {
		panic(err)
	}

	fmt.Println(esType)
	fmt.Println(dsl)
}
