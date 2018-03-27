package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	// get current username via env or command arg
	var username string
	flag.StringVar(&username, "u", os.Getenv("USER"), "specify your username to find your ssh config file")
	flag.Parse()

	file, err := os.Open("/Users/" + username + "/.ssh/config")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	str := string(buf)

	alias := regexp.MustCompile(`HOST\s(.*)\s`)
	ip := regexp.MustCompile(`HostName\s(.*)\s`)

	aliases := alias.FindAllString(str, -1)
	ips := ip.FindAllString(str, -1)
	for i := 0; i < len(ips); i++ {
		name := strings.Split(aliases[i], " ")[1]
		fmt.Printf(" \033[0;32m%-15s\033[0m%s", name[:len(name)-1], strings.Split(ips[i], " ")[1])
	}
}
