package client

// generateCoups generates coups for a given state and a given race
// While a player _can_ make multiple moves within a coup, for now this function only
// returns individual moves.
func generateCoups(s state, race race) []Coup {
	coups := []Coup{}

	for coord, cell := range s.grid {
		if cell.race == race {
			moves := generateMovesFromCell(s, coord)
			for _, move := range moves {
				coups = append(coups, Coup{move})
			}
		}
	}
	// TODO: generate more coups
	return coups
}

// transformation represents a coordinate transformation
// Ex: applying transform{-1, 1} to Coordinates{10, 11} would return Coordinates{9, 12}
type transformation struct {
	X int8
	Y int8
}

func transform(s state, c Coordinates, t transformation) (res Coordinates, ok bool) {
	if s.width == 0 {
		return c, false
	}

	xRes := uint8(int8(c.X) + t.X)
	if (xRes < 0) || (xRes >= s.width) {
		return c, false
	}

	yRes := uint8(int8(c.Y) + t.Y)
	if (yRes < 0) || (yRes >= s.height) {
		return c, false
	}

	return Coordinates{xRes, yRes}, true
}

func generateMovesFromCell(s state, source Coordinates) []Move {
	// Simplification: as long as we are doing one move per turn, we are better
	// off moving all units and not a subset.
	// This is not true anymore if we do multiple moves per turn, as keeping
	// some units in the source cell allows us to do another attack (on another
	// cell) in the same turn.

	moves := make([]Move, 0, 8)

	// TODO: can be optimized
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
			continue
		}
		moves = append(moves, Move{Start: source, N: s.grid[source].count, End: target})
	}
	return moves
}

// scoreState is the heuristic for our IA
func scoreState(potSta potentialState, race race) float64 {

	// Apply the change on the state

	h := 0.
	for _, cell := range potSta.s.grid {
		if cell.race == race {
			h += float64(cell.count)
		} else if cell.race == race.opponent() {
			h -= float64(cell.count)
		}
	}

	return h * potSta.probability
}
