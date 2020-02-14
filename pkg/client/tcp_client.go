package client

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/langorou/langorou/pkg/utils"
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

type ServerCmd int

const (
	_ ServerCmd = iota
	SET
	HUM
	HME
	MAP
	UPD
	END
	BYE
)

func (cmd ServerCmd) String() string {
	return [...]string{"UNKNOWN", "SET", "HUM", "HME", "MAP", "UPD", "END", "BYE"}[cmd]
}

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
func (c *TCPClient) SendName(name string) error {
	isASCII := utils.IsASCII(name)
	if !isASCII {
		return fmt.Errorf("%s is not a valid ASCII name", name)
	}

	t := len(name)
	if t == 0 || t > 255 {
		return fmt.Errorf("invalid name '%s', please use a short ASCII name", name)
	}

	msg := make([]byte, 3+1+t)
	fmt.Println(3 + 1 + t)
	copy(msg[:3], []byte("NME"))
	msg[3] = byte(uint8(t))
	copy(msg[4:], []byte(name))

	_, err := c.conn.Write(msg)

	return err
}

// SendMove to the server
func (c *TCPClient) SendMove(moves []Move) error {

	msg := make([]byte, 3+1+5*n)

	copy(msg[:3], []byte("MOV"))
	msg[3] = byte(n)

	for i := 0; i < int(n); i++ {
		msg[4+5*i] = byte(moves[i].Start.X)
		msg[4+5*i+1] = byte(moves[i].Start.Y)
		msg[4+5*i+2] = byte(moves[i].N)
		msg[4+5*i+3] = byte(moves[i].End.X)
		msg[4+5*i+4] = byte(moves[i].End.Y)
	}

	_, err := c.conn.Write(msg)

	return err
}

// ReceiveMsg from the server and parse it
func (c TCPClient) ReceiveMsg() (ServerCmd, error) {
	reader := bufio.NewReader(c.conn)
	buf := make([]byte, 9)
	if _, err := io.ReadFull(reader, buf[:3]); err != nil { // Read len(buf) chars
		return 0, err
	}

	command := string(buf[:3])
	switch command {
	case "SET":
		if _, err := io.ReadFull(reader, buf[:2]); err != nil {
			return SET, err
		}
		n := uint8(buf[0])
		m := uint8(buf[1])
		// DO something with it
		_, _ = n, m
		return SET, nil

	case "HUM":
		if _, err := io.ReadFull(reader, buf[:1]); err != nil {
			return HUM, err
		}
		n := uint8(buf[0])
		coords := make([]Coordinates, n)
		for i := 0; i < int(n); i++ {
			if _, err := io.ReadFull(reader, buf[:2]); err != nil {
				return HUM, err
			}
			coords[i] = Coordinates{
				X: uint8(buf[0]),
				Y: uint8(buf[1]),
			}
		}
		// DO something with it
		return HUM, nil

	case "HME":
		if _, err := io.ReadFull(reader, buf[:2]); err != nil {
			return HME, err
		}
		x := uint8(buf[0])
		y := uint8(buf[1])
		// DO something with it
		_, _ = x, y
		return HME, nil

	case "UPD":
		if _, err := io.ReadFull(reader, buf[:1]); err != nil {
			return UPD, err
		}
		n := uint8(buf[0])
		changes := make([]Changes, n)
		for i := 0; i < int(n); i++ {
			if _, err := io.ReadFull(reader, buf[:5]); err != nil {
				return UPD, err
			}
			changes[i] = Changes{
				X:          uint8(buf[0]),
				Y:          uint8(buf[1]),
				Humans:     uint(buf[2]),
				Vampires:   uint(buf[3]),
				Werewolves: uint(buf[4]),
			}
		}
		// DO something with it
		return UPD, nil

	case "MAP":
		if _, err := io.ReadFull(reader, buf[:1]); err != nil {
			return MAP, err
		}
		n := uint8(buf[0])
		changes := make([]Changes, n)
		for i := 0; i < int(n); i++ {
			if _, err := io.ReadFull(reader, buf[:5]); err != nil {
				return MAP, err
			}
			changes[i] = Changes{
				X:          uint8(buf[0]),
				Y:          uint8(buf[1]),
				Humans:     uint(buf[2]),
				Vampires:   uint(buf[3]),
				Werewolves: uint(buf[4]),
			}
		}

		// DO something with it

		return MAP, nil

	case "END":
		// Next Game
		return END, nil

	case "BYE":
		// Server stop
		return BYE, nil

	default:
		return 0, fmt.Errorf("invalid command from server : %s", command)
	}

}

// ReceiveSpecificCommand returns an error if the command is not as expected
func (c TCPClient) ReceiveSpecificCommand(assertCmd ServerCmd) error {
	command, err := c.ReceiveMsg()
	if err != nil {
		return err
	}
	if command != assertCmd {
		return fmt.Errorf("should have received %s but got %s instead", assertCmd, command)
	}
	return nil
}
