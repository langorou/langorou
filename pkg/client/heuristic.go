package client

type PotentialState struct {
	s    state
	prob float64
}

// generateMoves generates all possible *single* moves for a given state
// While a player _can_ make multiple moves, for now this function only
// returns individual moves.
func generateMoves(s state) []Move {
	moves := []Move{}

	for x, line := range s {
		for y, cell := range line {
			if cell.race == Ally {
				moves = append(moves, generateMovesFromCell(s, Coordinates{uint8(x), uint8(y)})...)
			}
		}
	}
	return moves
}

// transformation represents a coordinate transformation
// Ex: applying transform{-1, 1} to Coordinates{10, 11} would return Coordinates{9, 12}
type transformation struct {
	X int8
	Y int8
}

func transform(s state, c Coordinates, t transformation) (res *Coordinates, ok bool) {
	xLen := int8(len(s))
	if xLen == 0 {
		return nil, false
	}
	yLen := int8(len(s[0]))

	xRes := int8(c.X) + t.X
	if (xRes < 0) || (xRes >= xLen) {
		return nil, false
	}

	yRes := int8(c.Y) + t.Y
	if (yRes < 0) || (yRes >= yLen) {
		return nil, false
	}

	return &Coordinates{uint8(xRes), uint8(yRes)}, true
}

func generateMovesFromCell(s state, source Coordinates) []Move {
	// Simplification: as long as we are doing one move per turn, we are better
	// off moving all units and not a subset.
	// This is not true anymore if we do multiple moves per turn, as keeping
	// some units in the source cell allows us to do another attack (on another
	// cell) in the same turn.

	moves := []Move{}

	transforms := []transformation{
		{0, -1},  // left
		{0, 1},   // right
		{1, 0},   // top
		{-1, 0},  // bottom
		{1, -1},  // top-left
		{1, 1},   // top-right
		{-1, -1}, // bottom-left
		{-1, 1},  // bottom-right
	}

	for _, t := range transforms {
		target, ok := transform(s, source, t)
		if !ok {
			// move is out of bounds
			continue
		}
		moves = append(moves, Move{Start: source, N: s[source.X][source.Y].count, End: *target})
	}
	return moves
}

// applyMov to a state
func evaluateMoveOutcomes(s state, race race, coup Coup) []PotentialState {
	// we assume the mov is "safe"/correct

	// TODO: implement
	return []PotentialState{{s, 1}}
}

func scoreState(potSta PotentialState) float64 {

	// Apply the change on the state

	h := 0.
	for _, row := range potSta.s {
		for _, c := range row {
			switch c.race {
			case Empty:
				// nothing
			case Neutral:
				// nothing
			case Ally:
				h += float64(c.count)
			case Enemy:
				h -= float64(c.count)
			}
		}
	}

	return h * potSta.prob
}
