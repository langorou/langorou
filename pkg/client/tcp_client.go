package client

import (
	"fmt"
	"net"
)

type TcpClient struct {
	conn *net.Conn
}

// NewTcpClient creates a new TCP client
func NewTcpClient(ipAddr string, port string) (*TcpClient, error) {

	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipAddr)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ipAddr, port))
	if err != nil {
		return nil, err
	}

	return &TcpClient{
		conn: &conn,
	}, nil
}
