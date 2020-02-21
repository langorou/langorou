package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/langorou/langorou/pkg/utils"
)

// A good ressource
// https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/socket/tcp_sockets.html

// About net.Conn vs net.TCPConn :
// > This interface (net.Conn) has primary methods ReadFrom and WriteTo to handle packet reads and writes.
// > The Go net package recommends using these interface types rather than the concrete ones (TCPConn or UDPConn).
// > But by using them, you lose specific methods such as SetKeepAlive or TCPConn and SetReadBuffer of UDPConn, unless you do a type cast. It is your choice.

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

// TCPClient handles the connection to the server, and also encapsulate the game
type TCPClient struct {
	conn         net.Conn
	ourRaceCoord Coordinates
	isWerewolf   bool // We assume we're a vampire
	game *Game
}

// NewTCPClient creates a new TCP client
func NewTCPClient(addr string, name string, ia IA) (TCPClient, error) {

	// We might need to use net.DialTCP with https://golang.org/pkg/net/#TCPConn.SetKeepAlive
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return TCPClient{}, err
	}

	return TCPClient{
		conn: conn,
		game: NewGame(name, ia),
	}, nil
}

// SendName to the server
func (c *TCPClient) SendName() error {
	name := c.game.Nme()
	isASCII := utils.IsASCII(name)
	if !isASCII {
		return fmt.Errorf("%s is not a valid ASCII name", name)
	}

	t := len(name)
	if t == 0 || t > 255 {
		return fmt.Errorf("invalid name '%s', please use a short ASCII name", name)
	}

	msg := make([]byte, 3+1+t)
	copy(msg[:3], "NME")
	msg[3] = uint8(t)
	copy(msg[4:], name)

	_, err := c.conn.Write(msg)

	return err
}

// SendMove to the server
func (c *TCPClient) SendMove(moves []Move) error {

	n := len(moves)
	msg := make([]byte, 3+1+5*n)

	copy(msg[:3], "MOV")
	msg[3] = uint8(n)

	for i := 0; i < n; i++ {
		msg[4+5*i] = moves[i].Start.X
		msg[4+5*i+1] = moves[i].Start.Y
		msg[4+5*i+2] = moves[i].N
		msg[4+5*i+3] = moves[i].End.X
		msg[4+5*i+4] = moves[i].End.Y
	}

	_, err := c.conn.Write(msg)

	return err
}

// ReceiveMsg from the server and parse it
func (c *TCPClient) ReceiveMsg() (ServerCmd, error) {
	reader := bufio.NewReader(c.conn)
	buf := make([]byte, 5)                                  // we read at max 5 consecutive bytes
	if _, err := io.ReadFull(reader, buf[:3]); err != nil { // Read len(buf) chars, here 3 bytes
		return 0, err
	}

	command := string(buf[:3])
	switch command {
	case "SET":
		if _, err := io.ReadFull(reader, buf[:2]); err != nil {
			return SET, err
		}
		n := buf[0]
		m := buf[1]
		// DO something with it
		_, _ = n, m
		return SET, nil

	case "HUM":
		if _, err := io.ReadFull(reader, buf[:1]); err != nil {
			return HUM, err
		}
		n := buf[0]
		coords := make([]Coordinates, n)
		for i := 0; i < int(n); i++ {
			if _, err := io.ReadFull(reader, buf[:2]); err != nil {
				return HUM, err
			}
			coords[i] = Coordinates{
				X: buf[0],
				Y: buf[1],
			}
		}
		log.Printf("%s: Received %d positions of humans", command, len(coords))
		c.game.Hum(coords)
		return HUM, nil

	case "HME":
		if _, err := io.ReadFull(reader, buf[:2]); err != nil {
			return HME, err
		}
		x := buf[0]
		y := buf[1]

		c.ourRaceCoord = Coordinates{X: x, Y: y}
		log.Printf("%s: Received our race coordinates: %+v", command, c.ourRaceCoord)

		return HME, nil

	case "UPD":
		if _, err := io.ReadFull(reader, buf[:1]); err != nil {
			return UPD, err
		}
		n := buf[0]
		changes := make([]Changes, n)
		for i := 0; i < int(n); i++ {
			if _, err := io.ReadFull(reader, buf[:5]); err != nil {
				return UPD, err
			}

			if c.isWerewolf {
				changes[i] = Changes{
					Coords: Coordinates{
						X: buf[0],
						Y: buf[1],
					},
					Neutral: buf[2],
					Ally:    buf[4], // buf[4] represents the number of Werewolves
					Enemy:   buf[3], // buf[3] represents the number of Vampires
				}
			} else {
				changes[i] = Changes{
					Coords: Coordinates{
						X: buf[0],
						Y: buf[1],
					},
					Neutral: buf[2],
					Ally:    buf[3], // buf[3] represents the number of Vampires
					Enemy:   buf[4], // buf[4] represents the number of Werewolves
				}
			}
		}
		// DO something with it
		return UPD, nil

	case "MAP":
		if _, err := io.ReadFull(reader, buf[:1]); err != nil {
			return MAP, err
		}
		n := buf[0]
		changes := make([]Changes, n)

		flip := false // If we see that our start position is one of werewolf, we need to flip ally and enemy
		for i := 0; i < int(n); i++ {
			if _, err := io.ReadFull(reader, buf[:5]); err != nil {
				return MAP, err
			}

			changes[i] = Changes{
				Coords: Coordinates{
					X: buf[0],
					Y: buf[1],
				},
				Neutral: buf[2],
				Ally:    buf[3], // Vampire
				Enemy:   buf[4], // Werewolf
			}

			// Set to true if and only if we're actually werewolves
			flip = flip || (changes[i].Coords == c.ourRaceCoord && changes[i].Enemy > 0)

		}

		if flip {
			for i, c := range changes {
				c.Ally, c.Enemy = c.Enemy, c.Ally
				changes[i] = c
			}

			c.isWerewolf = true
		}

		// DO something with it

		return MAP, nil

	case "END":
		// Next Game

		// we reset some variables
		c.isWerewolf = false
		c.ourRaceCoord = Coordinates{}

		return END, nil

	case "BYE":
		// Server stop
		return BYE, nil

	default:
		return 0, fmt.Errorf("invalid command from server : %s", command)
	}

}

// ReceiveSpecificCommand returns an error if the command is not as expected
func (c *TCPClient) ReceiveSpecificCommand(assertCmd ServerCmd) error {
	command, err := c.ReceiveMsg()
	if err != nil {
		return err
	}
	if command != assertCmd {
		return fmt.Errorf("should have received %s but got %s instead", assertCmd, command)
	}
	return nil
}

// Init with the name
func (c *TCPClient) Init() error {
	// Send name
	err := c.SendName()
	if err != nil {
		return err
	}

	// Receive SET
	err = c.ReceiveSpecificCommand(SET)
	if err != nil {
		return err
	}

	// Receive HUM
	err = c.ReceiveSpecificCommand(HUM)
	if err != nil {
		return err
	}

	// Receive HME
	err = c.ReceiveSpecificCommand(HME)
	if err != nil {
		return err
	}

	// Receive MAP
	err = c.ReceiveSpecificCommand(MAP)
	if err != nil {
		return err
	}

	return nil
}
