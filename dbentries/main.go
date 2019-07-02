package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func main() {
	dbpath := flag.String("dbpath", "./", "please input the leveldb directory")
	flag.Parse()

	db, err := leveldb.OpenFile(*dbpath, nil)
	if err != nil {
		panic(err)
	}

	it := db.NewIterator(nil, &opt.ReadOptions{
		DontFillCache: true,
	})

	count := 0
	for it.Next() {
		count++
	}
	fmt.Printf("Total database entres: %d(%s)\n", count, *dbpath)

	it.Release()
	err = db.Close()
	if err != nil {
		fmt.Printf("close database error: %s\n", err)
		os.Exit(1)
	}
}

