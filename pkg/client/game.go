package client

import "log"

// Game implements the Client interface using a TCP server
type Game struct {
	state
	playerName string
}

// NewGame creates a new TCP client
func NewGame() (*Game, error) {
	return &Game{}, nil
}

// Nme defines the player name
func (g *Game) Nme(playerName string) string {
	g.playerName = playerName
	return playerName
}

func (T Game) Mov() []Move {
	panic("implement me")
}

// Set initialize an empty grid in the state
func (g *Game) Set(n uint8, m uint8) {
	g.state.grid = make([][]cell, n)
	for i := range g.state.grid {
		g.state.grid[i] = make([]cell, m)
	}
}

func (g *Game) Hum(coords []Coordinates) {
	for _, pos := range coords {
		g.state.grid[pos.Y][pos.X].race = Neutral
		//TODO check map and human order
		g.state.grid[pos.Y][pos.X].count = 0
	}
}

//Upd updates the state of the game
func (g *Game) Upd(changes []Changes) {

	for _, cha := range changes {
		if cha.Neutral > 0 && cha.Ally == 0 && cha.Enemy == 0 {
			g.state.grid[cha.Coords.Y][cha.Coords.X].count = cha.Neutral
			g.state.grid[cha.Coords.Y][cha.Coords.X].race = Neutral
		} else if cha.Neutral == 0 && cha.Ally > 0 && cha.Enemy == 0 {
			g.state.grid[cha.Coords.Y][cha.Coords.X].count = cha.Ally
			g.state.grid[cha.Coords.Y][cha.Coords.X].race = Ally
		} else if cha.Neutral == 0 && cha.Ally == 0 && cha.Enemy > 0 {
			g.state.grid[cha.Coords.Y][cha.Coords.X].count = cha.Enemy
			g.state.grid[cha.Coords.Y][cha.Coords.X].race = Enemy
		} else if cha.Neutral == 0 && cha.Ally == 0 && cha.Enemy == 0 {
			g.state.grid[cha.Coords.Y][cha.Coords.X].count = 0
			g.state.grid[cha.Coords.Y][cha.Coords.X].race = Empty
		} else {
			log.Printf("impossible case: only one race on a cell")
		}
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
