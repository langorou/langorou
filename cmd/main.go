package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/langorou/langorou/pkg/client"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Printf("please provide IP address and port\n")
		os.Exit(1)
	}

	ip := net.ParseIP(args[0])
	if ip == nil {
		fmt.Printf("invalid IP address: %s\n", args[0])
		os.Exit(1)
	}

	_, err := strconv.ParseUint(args[1], 10, 16) // 0 <= port <= 65535
	if err != nil {
		fmt.Printf("invalid port %s, should be between 0 and 65535\n", args[1])
		os.Exit(1)
	}

	addr := net.JoinHostPort(ip.String(), args[1])
	fmt.Printf("connecting to %s\n", addr)

	os.Exit(0)
	_, err = client.NewTCPClient(ip, args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
