package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	args := os.Args

	if len(args) <= 1 {
		fmt.Printf("timestamp: %d\n", time.Now().Unix())
		return
	}

	timestamp, err := strconv.Atoi(args[1])
	// timestamp to time string
	if err == nil {
		ts := time.Unix(int64(timestamp), 0)
		fmt.Printf("time: %s\n", ts.Format("2006-01-02 15:04:05"))
		return
	}

	// time to timestamp
	t, err := time.Parse("2006-01-02 15:04:05", args[1])
	if err == nil {
		fmt.Printf("timestamp: %d\n", t.Unix())
		return
	}

	fmt.Println("invalid parameter!")
}
