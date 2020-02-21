package client

import (
	"fmt"
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

type ServerCmd string

const (
	UNKNOWN ServerCmd = "UNKNOWN"
	SET               = "SET"
	HUM               = "HUM"
	HME               = "HME"
	MAP               = "MAP"
	UPD               = "UPD"
	END               = "END"
	BYE               = "BYE"
)

// TCPClient handles the connection to the server, and also encapsulate the game
type TCPClient struct {
	conn         net.Conn
	ourRaceCoord Coordinates
	isWerewolf   bool // We assume we're a vampire
	game         *Game
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

	msg := make([]byte, 3+1+5*len(moves))

	copy(msg[:3], "MOV")
	msg[3] = uint8(len(moves))

	log.Printf("===")
	log.Printf("Sending %d moves:", len(moves))
	for i, move := range moves {
		log.Printf("Move: %+v", move)
		msg[4+5*i] = move.Start.X
		msg[4+5*i+1] = move.Start.Y
		msg[4+5*i+2] = move.N
		msg[4+5*i+3] = move.End.X
		msg[4+5*i+4] = move.End.Y
	}
	log.Printf("===")

	_, err := c.conn.Write(msg)

	return err
}

// ReceiveMsg from the server and parse it
func (c *TCPClient) ReceiveMsg() (ServerCmd, error) {
	buf := make([]byte, 5)                          // we read at max 5 consecutive bytes
	if _, err := c.conn.Read(buf[:3]); err != nil { // Read len(buf) chars, here 3 bytes
		return UNKNOWN, err
	}

	command := string(buf[:3])
	log.Printf("Received command: %s", command)
	switch command {
	case "SET":
		if _, err := c.conn.Read(buf[:2]); err != nil {
			return SET, err
		}
		n := buf[0]
		m := buf[1]

		log.Printf("%s: set the map size to (%d, %d)", command, n, m)

		c.game.Set(n, m)

		return SET, nil

	case "HUM":
		if _, err := c.conn.Read(buf[:1]); err != nil {
			return HUM, err
		}
		n := buf[0]
		coords := make([]Coordinates, n)
		for i := 0; i < int(n); i++ {
			if _, err := c.conn.Read(buf[:2]); err != nil {
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
		if _, err := c.conn.Read(buf[:2]); err != nil {
			return HME, err
		}
		x := buf[0]
		y := buf[1]

		c.ourRaceCoord = Coordinates{X: x, Y: y}
		log.Printf("%s: Received our race coordinates: %+v", command, c.ourRaceCoord)

		return HME, nil

	case "UPD":
		if _, err := c.conn.Read(buf[:1]); err != nil {
			return UPD, err
		}
		n := buf[0]
		changes := make([]Changes, n)
		for i := 0; i < int(n); i++ {
			if _, err := c.conn.Read(buf[:5]); err != nil {
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
		log.Printf("%s: received %d changes", command, n)
		c.game.Upd(changes)
		return UPD, nil

	case "MAP":
		if _, err := c.conn.Read(buf[:1]); err != nil {
			return MAP, err
		}
		n := buf[0]
		changes := make([]Changes, n)

		flip := false // If we see that our start position is one of werewolf, we need to flip ally and enemy
		for i := 0; i < int(n); i++ {
			if _, err := c.conn.Read(buf[:5]); err != nil {
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
		log.Printf("%s: received %d changes", command, n)

		if flip {
			log.Printf("%s: we are werewolves", command)
			for i, c := range changes {
				c.Ally, c.Enemy = c.Enemy, c.Ally
				changes[i] = c
			}

			c.isWerewolf = true
		}

		c.game.Map(changes)

		return MAP, nil

	case "END":
		// Next Game
		log.Printf("%s: end of the game", command)
		// we reset some variables
		c.isWerewolf = false
		c.ourRaceCoord = Coordinates{}

		return END, c.game.End()

	case "BYE":
		// Server stop
		log.Printf("%s: server said bye", command)
		return BYE, nil

	default:
		return UNKNOWN, fmt.Errorf("invalid command from server : %s", command)
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

// Start with the name
func (c *TCPClient) Start() error {
	// TODO: it's possible to receive BYE here, if we restarted a game

	if err := c.init(); err != nil {
		return fmt.Errorf("an error occurred during init: %s", err)
	}
	log.Print("Successfully init !")

	log.Printf("TODO")

	return nil
}

func (c *TCPClient) init() error {
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

	for {
		cmd, err := c.ReceiveMsg()
		if err != nil {
			return err
		}

		switch cmd {
		case UPD:
			if err = c.SendMove(c.game.Mov()); err != nil {
				return err
			}
		case BYE:
			log.Printf("Received BYE, stopping the client...")
			return nil
		case END:
			log.Printf("Received END, getting ready for the next game...")
			return c.Start()
		default:
			return fmt.Errorf("received unexpected command: %s", cmd)
		}
	}
}
