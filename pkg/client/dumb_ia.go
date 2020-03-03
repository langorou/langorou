package client

import "time"

type DumbIA struct{}

var _ IA = &DumbIA{}

func NewDumbIA() *DumbIA {
	return &DumbIA{}
}

func (dia *DumbIA) Play(state state) Coup {
	// Simulate computation
	time.Sleep(time.Second)

	for coord, cell := range state.grid {
		x := coord.X
		y := coord.Y
		if cell.race == Ally && cell.count > 0 {

			endY := y + 1
			endX := x + 1
			if endY == state.height {
				endY -= 2
			}
			if endX == state.width {
				endX -= 2
			}

			return []Move{
				{
					Start: coord,
					// Everyone moves !
					N:   cell.count,
					End: Coordinates{X: endX, Y: endY},
				},
			}
		}
	}

	panic("SHOULD NOT HAPPEN ! DUMB IA did not find any valid move")
}
