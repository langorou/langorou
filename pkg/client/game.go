package client

import (
	"fmt"
	"github.com/langorou/langorou/pkg/client/model"
)

// Game implements the Client interface using a TCP server
type Game struct {
	state      model.State
	playerName string
	ia         IA
}

// NewGame creates a new TCP client
func NewGame(name string, ia IA) *Game {
	return &Game{playerName: name, ia: ia}
}

// Nme defines the player name
func (g *Game) Nme() string {
	return g.playerName
}

func (g *Game) Mov() []model.Move {
	return g.ia.Play(g.state)
}

// Set initialize an empty grid in the state
func (g *Game) Set(n uint8, m uint8) {
	g.state = model.NewState(n, m)
}

func (g *Game) Hum(coords []model.Coordinates) {
	for _, pos := range coords {
		g.state.SetCell(pos, model.Neutral, 0)
	}
}

//Upd updates the state of the game
func (g *Game) Upd(changes []model.Changes) {

	for _, change := range changes {
		if change.Neutral > 0 && change.Ally == 0 && change.Enemy == 0 {
			// Neutral
			g.state.SetCell(change.Coords, model.Neutral, change.Neutral)
		} else if change.Neutral == 0 && change.Ally > 0 && change.Enemy == 0 {
			// Ally
			g.state.SetCell(change.Coords, model.Ally, change.Ally)
		} else if change.Neutral == 0 && change.Ally == 0 && change.Enemy > 0 {
			// Enemy
			g.state.SetCell(change.Coords, model.Enemy, change.Enemy)
		} else if change.Neutral == 0 && change.Ally == 0 && change.Enemy == 0 {
			// Empty Cell
			g.state.EmptyCell(change.Coords)
		} else {
			// Should not happen !
			panic(fmt.Sprintf("impossible change, maximum one race per cell: %+v", change))
		}
	}
}

// Map is the same as Upd but is called only once at the beginning
func (g *Game) Map(changes []model.Changes) {
	g.Upd(changes)
}

// End delete the state of the game
func (g *Game) End() error {
	g = &Game{}
	return nil
}

var _ = &Game{}
