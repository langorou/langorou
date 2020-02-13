package client

import (
	"fmt"
	"net"
)

// TCPClient connects to the game server
type TCPClient struct {
	conn *net.Conn
}

// NewTCPClient creates a new TCP client
func NewTCPClient(ipAddr net.IP, port string) (*TCPClient, error) {

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ipAddr, port))
	if err != nil {
		return nil, err
	}

	return &TCPClient{
		conn: &conn,
	}, nil
}
