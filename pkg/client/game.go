package client

import "log"

// Game implements the Client interface using a TCP server
type Game struct {
	state      state
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

func (g *Game) Mov() []Move {
	return g.ia.Play(g.state)
}

// Set initialize an empty grid in the state
func (g *Game) Set(n uint8, m uint8) {
	g.state = state{
		grid:   map[Coordinates]cell{},
		height: n,
		width:  m,
	}
}

func (g *Game) Hum(coords []Coordinates) {
	for _, pos := range coords {
		cell := g.state.grid[pos]
		cell.race = Neutral
		cell.count = 0
		g.state.grid[pos] = cell
	}
}

//Upd updates the state of the game
func (g *Game) Upd(changes []Changes) {

	for _, change := range changes {
		cell := g.state.grid[change.Coords]
		if change.Neutral > 0 && change.Ally == 0 && change.Enemy == 0 {
			cell.count = change.Neutral
			cell.race = Neutral
		} else if change.Neutral == 0 && change.Ally > 0 && change.Enemy == 0 {
			cell.count = change.Ally
			cell.race = Ally
		} else if change.Neutral == 0 && change.Ally == 0 && change.Enemy > 0 {
			cell.count = change.Enemy
			cell.race = Enemy
		} else {
			log.Printf("impossible change, maximum one race per cell: %+v", change)
			delete(g.state.grid, change.Coords)
		}
		g.state.grid[change.Coords] = cell
	}
}

// Map is the same as Upd but is called only once at the beginning
func (g *Game) Map(changes []Changes) {
	g.Upd(changes)
}

// End delete the state of the game
func (g *Game) End() error {
	g = &Game{}
	return nil
}

var _ = &Game{}
