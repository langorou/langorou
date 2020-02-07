package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("please provide IP address and port")
		os.Exit(1)
	}
	ipAddr := args[0]
	port := args[1]
	fmt.Printf("ip %s port %s\n", ipAddr, port)
	os.Exit(0)
}
