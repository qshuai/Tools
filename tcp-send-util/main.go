package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	host := flag.String("host", "127.0.0.1", "foreign address")
	port := flag.Int("port", 80, "foreign host listen port")
	length := flag.Int("length", 64, "how many bytes to send")
	sleep := flag.Int64("sleep", 0, "how many milliseconds to hang between connecting and sending")
	sleepAfterWrite := flag.Int64("sleep-after-write", 3000, "how many milliseconds to wait before exiting")
	writeTimeout := flag.Int64("write-timeout", 1000, "how many milliseconds to wait before returning writing successful")
	flag.Parse()

	// 参数校验
	if len(*host) <= 0 {
		fmt.Println("Empty host not accepted")
		os.Exit(1)
	}
	if *port <= 0 {
		fmt.Println("Invalid tcp port")
		os.Exit(1)
	}
	if *length <= 0 {
		fmt.Println("Specified length should more than 0")
		os.Exit(1)
	}

	ip := net.ParseIP(*host)
	if ip == nil {
		fmt.Println("Please input a valid ip addr(ipv4/ipv6)")
		os.Exit(1)
	}

	addr := net.TCPAddr{
		IP:   ip,
		Port: *port,
	}
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		fmt.Println("Connection error:", err)
		os.Exit(1)
	}

	// whether hanging or not
	if *sleep > 0 {
		time.Sleep(time.Duration(*sleep) * time.Millisecond)
	}

	// setsockopt after sleeping
	err = conn.SetWriteDeadline(time.Now().Add(time.Duration(*writeTimeout) * time.Millisecond))
	if err != nil {
		fmt.Println("Set write timeout socketopt err:", err)
		os.Exit(1)
	}

	fmt.Println("Start to send...")
	n, err := conn.Write(make([]byte, *length))
	if err != nil {
		fmt.Println("Write to socket err:", err)
		os.Exit(1)
	}

	fmt.Println("Send Successful, data length:", n)

	if *sleepAfterWrite > 0 {
		time.Sleep(time.Duration(*sleepAfterWrite) * time.Millisecond)
	}
}
