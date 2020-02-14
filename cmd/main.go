package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/langorou/langorou/pkg/client"
)

func failIf(err error, msg string) {
	if err != nil {
		log.Fatalf("error %s: %v", msg, err)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Printf("please provide IP address and port\n")
		os.Exit(1)
	}

	// We might use localhost...
	// ip := net.ParseIP(args[0])
	// if ip == nil {
	// 	fmt.Printf("invalid IP address: %s\n", args[0])
	// 	os.Exit(1)
	// }

	_, err := strconv.ParseUint(args[1], 10, 16) // 0 <= port <= 65535
	failIf(err, fmt.Sprintf("invalid port %s, should be between 0 and 65535\n", args[1]))

	addr := net.JoinHostPort(args[0], args[1])
	fmt.Printf("connecting to %s\n", addr)

	c, err := client.NewTCPClient(addr)
	failIf(err, "")

	// Send name
	err = c.SendName("langorou")
	failIf(err, "")

	// Receive SET
	err = c.ReceiveSpecificCommand(client.SET)
	failIf(err, "")

	// Receive HUM
	err = c.ReceiveSpecificCommand(client.HUM)
	failIf(err, "")

	// Receive HME
	err = c.ReceiveSpecificCommand(client.HME)
	failIf(err, "")

	// Receive MAP
	err = c.ReceiveSpecificCommand(client.MAP)
	failIf(err, "")

	os.Exit(0)
}
