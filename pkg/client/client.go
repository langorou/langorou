package client

import (
	"github.com/langorou/langorou/pkg/client/model"
)

// Player represents the player -> server protocol
type Player interface {
	// Player -> Server
	Nme() (name string)
	Mov() (moves []model.Move)
}

// Server represents the server -> player protocol
type Server interface {
	// Server -> Player
	Set(n uint8, m uint8)
	Hum(coords []model.Coordinates)
	// Hme(x uint8, y uint8) handled by TCP Client
	Upd(changes []model.Changes)
	Map(changes []model.Changes)
	End() error
}

// Client implements the Player and Server interfaces
type Client interface {
	Player
	Server
}
