package client

type IA interface {
	Play(state state) []Move
}

type DumbIA struct {}

func NewDumbIA() *DumbIA {
	return &DumbIA{}
}

func (dia *DumbIA) Play(state state) []Move {
	for y, row := range state {
		for x, cell := range row {
			if cell.race == Ally && cell.count > 0 {

				endY := y + 1
				endX := x +1
				if y == len(state) {
					endY -= 2
				}
				if x == len(row) {
					endX -= 2
				}

				return []Move{
					{
						// TODO: figure out if this is correct for y
						Start: Coordinates{X: uint8(x), Y: uint8(y)},
						// Everyone moves !
						N: uint8(cell.count),
						End: Coordinates{X: uint8(endX), Y: uint8(endY)},
					},
				}
			}
		}
	}

	panic("SHOULD NOT HAPPEN ! DUMB IA did not find any valid move")
}