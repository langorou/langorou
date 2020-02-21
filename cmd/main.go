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
	namePtr := flag.String("name", "langorou", "name of the player")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Printf("please provide IP address and port\n")
		os.Exit(1)
	}

	_, err := strconv.ParseUint(args[1], 10, 16) // 0 <= port <= 65535
	failIf(err, fmt.Sprintf("invalid port %s, should be between 0 and 65535\n", args[1]))

	addr := net.JoinHostPort(args[0], args[1])
	log.Printf("connecting to %s with name: %s", addr, *namePtr)

	c, err := client.NewTCPClient(addr, *namePtr, client.NewDumbIA())
	failIf(err, "")

	failIf(c.Start(), "")

	os.Exit(0)
}
