package client

import "time"

type DumbIA struct{}

func NewDumbIA() *DumbIA {
	return &DumbIA{}
}

func (dia *DumbIA) Play(state state) Coup {
	// Simulate computation
	time.Sleep(time.Second)

	for y, row := range state {
		for x, cell := range row {
			if cell.race == Ally && cell.count > 0 {

				endY := y + 1
				endX := x + 1
				if endY == len(state) {
					endY -= 2
				}
				if endX == len(row) {
					endX -= 2
				}

				return []Move{
					{
						Start: Coordinates{X: uint8(x), Y: uint8(y)},
						// Everyone moves !
						N:   uint8(cell.count),
						End: Coordinates{X: uint8(endX), Y: uint8(endY)},
					},
				}
			}
		}
	}

	panic("SHOULD NOT HAPPEN ! DUMB IA did not find any valid move")
}
