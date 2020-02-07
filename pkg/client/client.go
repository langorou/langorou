package client

// Coordinates represents coordinates on the grid
type Coordinates struct {
	X uint8
	Y uint8
}

// Changes is an update of the board sent by the server
type Changes struct {
	X          uint8
	Y          uint8
	Humans     uint
	Vampires   uint
	Werewolves uint
}

// Move is an allowed move
type Move struct {
	Start Coordinates
	N     uint8
	End   Coordinates
}

// Player represents the player -> server protocol
type Player interface {
	// Player -> Server
	Nme(t uint8, name []byte) error
	Mov(n uint8, moves []Move) error
}

// Server represents the server -> player protocol
type Server interface {
	// Server -> Player
	Set(n uint8, m uint8) error
	Hum(n uint8, coords []Coordinates) error
	Hme(x uint8, y uint8) error
	Upd(n uint8, changes []Changes) error
	Map(n uint8, changes []Changes) error
	End() error
}

// Client implements the Player and Server interfaces
type Client interface {
	Player
	Server
}
