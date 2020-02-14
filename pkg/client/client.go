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
	Humans     uint8
	Vampires   uint8
	Werewolves uint8
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
	Nme() (name string)
	Mov() (moves []Move)
}

// Server represents the server -> player protocol
type Server interface {
	// Server -> Player
	Set(n uint8, m uint8)
	Hum(coords []Coordinates)
	Hme(x uint8, y uint8)
	Upd(changes []Changes)
	Map(changes []Changes)
	End() error
}

// Client implements the Player and Server interfaces
type Client interface {
	Player
	Server
}
