package client

import (
	"math"
)

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

// scoreNeutralBattle scores the issue of a battle between a monster and a neutral group
func scoreNeutralBattle(c1, c2 Coordinates, cell1, cell2 cell) float64 {
	proba := winProbability(cell1.count, cell2.count, true)
	distance := c1.Distance(c2)

	// probable gain of population
	probableGain := math.Max(
		0,
		proba*float64(cell1.count+cell2.count)-float64(cell1.count),
	)

	return probableGain / distance
}

// scoreMonsterBattle scores the issue of a battle between monsters
func scoreMonsterBattle(c1, c2 Coordinates, cell1, cell2 cell) (float64, float64) {
	distance := c1.Distance(c2)

	// p1 is for 1 attacks 2
	p1 := winProbability(cell1.count, cell2.count, false)
	// p2 is for 2 attacks 1
	p2 := winProbability(cell2.count, cell1.count, false)

	s1 := p1*float64(cell1.count+cell2.count) - float64(cell2.count)
	s2 := p2*float64(cell2.count+cell1.count) - float64(cell1.count)

	return s1 / distance, s2 / distance
}

// scoreState is the heuristic for our IA
func scoreState(potSta potentialState, ourRace race) float64 {

	// different counts participating in the heuristic
	counts := map[race]float64{Ally: 0, Enemy: 0}
	battleCounts := map[race]float64{Ally: 0, Enemy: 0}
	neutralBattleCounts := map[race]float64{Ally: 0, Enemy: 0}

	for c1, cell1 := range potSta.s.grid {
		if cell1.race == Neutral {
			continue
		}
		counts[cell1.race] += float64(cell1.count)

		// Loop to compute stats on the possible battle
		for c2, cell2 := range potSta.s.grid {
			if c1 == c2 || cell1.race == cell2.race {
				continue
			}

			if cell2.race == Neutral {
				// TODO: average here since we can count a battle multiple times
				neutralBattleCounts[cell1.race] += scoreNeutralBattle(c1, c2, cell1, cell2)
			} else if cell2.race == cell1.race.opponent() {
				// TODO: average here since we can count a battle multiple times
				g1, g2 := scoreMonsterBattle(c1, c2, cell1, cell2)
				battleCounts[cell1.race] += g1
				battleCounts[cell2.race] += g2
			}
		}
	}

	total := 0.

	// TODO: make those parameters of a heuristic struct and try to tweak them
	// TODO: distance power alpha instead of distance power 1
	const (
		countCoeff         = 2
		battleCoeff        = 0.2
		neutralBattleCoeff = 0.5
	)

	for _, heuristic := range []struct {
		coeff  float64
		scores map[race]float64
	}{
		{countCoeff, counts},
		{battleCoeff, battleCounts},
		{neutralBattleCoeff, neutralBattleCounts},
	} {
		score := 0.
		for race, count := range heuristic.scores {
			if race == ourRace {
				score += count
			} else if race.opponent() == ourRace {
				score -= count
			}
		}
		total += score * heuristic.coeff
	}

	return total * potSta.probability
}
