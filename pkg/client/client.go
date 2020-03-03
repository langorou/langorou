package client

// Coordinates represents coordinates on the grid
type Coordinates struct {
	X uint8
	Y uint8
}

// Changes is an update of the board sent by the server
type Changes struct {
	Coords  Coordinates
	Neutral uint8
	Ally    uint8
	Enemy   uint8
}

// Move is an allowed move
type Move struct {
	Start Coordinates
	N     uint8
	End   Coordinates
}

// Coup represents a list of moves/actions, it implements the sort.Interface to sort by target cells
type Coup []Move

func (coup Coup) Len() int {
	return len(coup)
}

func (coup Coup) Less(i, j int) bool {
	return coup[i].End.Y > coup[j].End.Y && coup[i].End.X > coup[j].End.X
}

func (coup Coup) Swap(i, j int) {
	coup[i], coup[j] = coup[j], coup[i]
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
	// Hme(x uint8, y uint8) handled by TCP Client
	Upd(changes []Changes)
	Map(changes []Changes)
	End() error
}

// Client implements the Player and Server interfaces
type Client interface {
	Player
	Server
}
