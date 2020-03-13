package client

import (
	"github.com/langorou/langorou/pkg/client/model"
	"math"
)

// generateCoups generates coups for a given state and a given race
// While a player _can_ make multiple moves within a coup, for now this function only
// returns individual moves.
func generateCoups(s model.State, race model.Race) []model.Coup {
	coups := []model.Coup{}

	for coord, cell := range s.Grid {
		if cell.Race == race {
			moves := generateMovesFromCell(s, coord)
			for _, move := range moves {
				coups = append(coups, model.Coup{move})
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

func transform(s model.State, c model.Coordinates, t transformation) (res model.Coordinates, ok bool) {
	if s.Width == 0 {
		return c, false
	}

	xRes := uint8(int8(c.X) + t.X)
	if (xRes < 0) || (xRes >= s.Width) {
		return c, false
	}

	yRes := uint8(int8(c.Y) + t.Y)
	if (yRes < 0) || (yRes >= s.Height) {
		return c, false
	}

	return model.Coordinates{X: xRes, Y: yRes}, true
}

func generateMovesFromCell(s model.State, source model.Coordinates) []model.Move {
	// Simplification: as long as we are doing one move per turn, we are better
	// off moving all units and not a subset.
	// This is not true anymore if we do multiple moves per turn, as keeping
	// some units in the source cell allows us to do another attack (on another
	// cell) in the same turn.

	moves := make([]model.Move, 0, 8)

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
		moves = append(moves, model.Move{Start: source, N: s.Grid[source].Count, End: target})
	}
	return moves
}

// scoreNeutralBattle scores the issue of a battle between a monster and a neutral group
func scoreNeutralBattle(c1, c2 model.Coordinates, cell1, cell2 model.Cell) float64 {
	proba := winProbability(cell1.Count, cell2.Count, true)
	distance := c1.Distance(c2)

	// probable gain of population
	probableGain := math.Max(
		0,
		proba*float64(cell1.Count+cell2.Count)-float64(cell1.Count),
	)

	return probableGain / distance
}

// scoreMonsterBattle scores the issue of a battle between monsters
func scoreMonsterBattle(c1, c2 model.Coordinates, cell1, cell2 model.Cell) (float64, float64) {
	distance := c1.Distance(c2)

	// p1 is for 1 attacks 2
	p1 := winProbability(cell1.Count, cell2.Count, false)
	// p2 is for 2 attacks 1
	p2 := winProbability(cell2.Count, cell1.Count, false)

	s1 := p1*float64(cell1.Count+cell2.Count) - float64(cell2.Count)
	s2 := p2*float64(cell2.Count+cell1.Count) - float64(cell1.Count)

	return s1 / distance, s2 / distance
}

type scoreCounter struct {
	ally  float64
	enemy float64
}

func (sc *scoreCounter) add(race model.Race, score float64) {
	switch race {
	case model.Ally:
		sc.ally += score
	case model.Enemy:
		sc.enemy += score
	}
}

// scoreState is the heuristic for our IA
func scoreState(s model.State) float64 {

	// different counts participating in the heuristic
	counts := scoreCounter{}
	battleCounts := scoreCounter{}
	neutralBattleCounts := scoreCounter{}

	for c1, cell1 := range s.Grid {
		if cell1.Race == model.Neutral {
			continue
		}
		counts.add(cell1.Race, float64(cell1.Count))

		// Loop to compute stats on the possible battle
		for c2, cell2 := range s.Grid {
			if c1 == c2 || cell1.Race == cell2.Race {
				continue
			}

			if cell2.Race == model.Neutral {
				// TODO: average here since we can count a battle multiple times
				neutralBattleCounts.add(cell1.Race, scoreNeutralBattle(c1, c2, cell1, cell2))
			} else if cell2.Race == cell1.Race.Opponent() {
				// TODO: average here since we can count a battle multiple times
				g1, g2 := scoreMonsterBattle(c1, c2, cell1, cell2)
				battleCounts.add(cell1.Race, g1)
				battleCounts.add(cell2.Race, g2)
			}
		}
	}

	total := 0.

	// TODO: make those parameters of a heuristic struct and try to tweak them
	// TODO: distance power alpha instead of distance power 1
	const (
		countCoeff         = 1
		battleCoeff        = 0.5
		neutralBattleCoeff = 0.5
		cumScoreCoeff      = 0.001
		winScore = 1e10
		loseOverWinRatio = 1.2
	)

	cumScore := s.CumulativeScore

	// Win and lose cases
	if counts.ally == 0 {
		return -(winScore*loseOverWinRatio) + cumScore
	} else if counts.enemy == 0 {
		return winScore + cumScore
	}

	for _, heuristic := range []struct {
		coeff  float64
		scores scoreCounter
	}{
		{countCoeff, counts},
		{battleCoeff, battleCounts},
		{neutralBattleCoeff, neutralBattleCounts},
	} {
		score := heuristic.scores.ally - heuristic.scores.enemy
		total += score * heuristic.coeff
	}

	return total + (cumScore * cumScoreCoeff)
}
