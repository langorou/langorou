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

func (g *Game) Hum(n uint8, coords []Coordinates) {
	for _, pos := range coords {
		g.state.grid[pos.Y][pos.X].race = Villager
		//TODO check map and human order
		g.state.grid[pos.Y][pos.X].count = 0
	}
}

func (T Game) Hme(x uint8, y uint8) error {
	panic("implement me")
}

//Upd updates the state of the game
func (g *Game) Upd(changes []Changes) {

	for _, cha := range changes {
		if cha.Humans > 0 && cha.Vampires == 0 && cha.Werewolves == 0 {
			g.state.grid[cha.Y][cha.X].count = cha.Humans
			g.state.grid[cha.Y][cha.X].race = Villager
		} else if cha.Humans == 0 && cha.Vampires > 0 && cha.Werewolves == 0 {
			g.state.grid[cha.Y][cha.X].count = cha.Vampires
			g.state.grid[cha.Y][cha.X].race = Vampire
		} else if cha.Humans == 0 && cha.Vampires == 0 && cha.Werewolves > 0 {
			g.state.grid[cha.Y][cha.X].count = cha.Werewolves
			g.state.grid[cha.Y][cha.X].race = Werewolf
		} else if cha.Humans == 0 && cha.Vampires == 0 && cha.Werewolves == 0 {
			g.state.grid[cha.Y][cha.X].count = 0
			g.state.grid[cha.Y][cha.X].race = Empty
		} else {
			log.Printf("impossible case: only one race on a cell")
		}
	}
}

//Map is the same as Upd but is called only once at the beginning
func (g *Game) Map(changes []Changes) {
	g.Upd(changes)
}

func (T Game) End() error {
	panic("implement me")
}

var _ = &Game{}
