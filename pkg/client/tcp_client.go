package client

import (
	"fmt"
	"net"
	"strconv"
)

// A good ressource
// https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/socket/tcp_sockets.html

// About net.Conn vs net.TCPConn :
// > This interface has primary methods ReadFrom and WriteTo to handle packet reads and writes.
// > The Go net package recommends using these interface types rather than the concrete ones.
// > But by using them, you lose specific methods such as SetKeepAlive or TCPConn and SetReadBuffer of UDPConn, unless you do a type cast. It is your choice.

// TCPClient connects to the game server

type ClientMsg int

const (
	NME ClientMsg = iota
	MOV
)

type ServerMsg int

const (
	SET ServerMsg = iota
	HUM
	HME
	MAP
	UPD
	END
	BYE
)

type TCPClient struct {
	conn net.Conn
}

// NewTCPClient creates a new TCP client
func NewTCPClient(addr string) (TCPClient, error) {

	// We might need to use net.DialTCP with https://golang.org/pkg/net/#TCPConn.SetKeepAlive
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return TCPClient{}, err
	}

	return TCPClient{
		conn: conn,
	}, nil
}

// SendName to the server
func (c TCPClient) SendName(name string) error {
	ASCIIName := strconv.QuoteToASCII(name)
	t := len(ASCIIName)
	if t == 0 || t > 255 {
		return fmt.Errorf("invalid name '%s', please use a short ASCII name", name)
	}
	_, err := c.conn.Write([]byte(fmt.Sprintf("NME%d%s", t, ASCIIName)))

	return err
}

// SendMove to the server
func (c TCPClient) SendMove(n uint8, moves []Move) error {
	strMsg := fmt.Sprintf("MOV%d", n)

	for i := 0; i < int(n); i++ {
		strMsg += fmt.Sprintf("%d%d%d%d%d", moves[i].Start.X, moves[i].Start.Y, moves[i].N, moves[i].End.X, moves[i].End.Y)
	}

	msg := []byte(strMsg)

	_, err := c.conn.Write(msg)

	return err
}
