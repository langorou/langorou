package client

// Game implements the Client interface using a TCP server
type Game struct {
	state
}

// NewGame creates a new TCP client
func NewGame() (*Game, error) {
	return &Game{}, nil
}

func (T Game) Nme() string {
	panic("implement me")
}

func (T Game) Mov() []Move {
	panic("implement me")
}

func (T Game) Set(n uint8, m uint8) error {
	panic("implement me")
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

func (T Game) Upd(n uint8, changes []Changes) error {
	panic("implement me")
}

func (T Game) Map(n uint8, changes []Changes) error {
	panic("implement me")
}

func (T Game) End() error {
	panic("implement me")
}

var _ = &Game{}
