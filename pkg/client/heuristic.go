package client

import (
	"fmt"
	"math"

	"github.com/langorou/langorou/pkg/client/model"
)

// HeuristicParameters defines the parameters used to compute the heuristic
type HeuristicParameters struct {
	// Counts is a coefficient for the counts of populations
	Counts float64
	// Battles is a coefficient for battles against opponents
	Battles float64
	// NeutralBattles is a coefficient for battles against neutrals
	NeutralBattles float64
	// CumScore is a coefficient for the cumulative score (used to give more priority to scores achieved with shorter paths)
	CumScore float64

	// WinScore given to a win
	WinScore float64

	// Ratio of LoseScore / WinScore, allows to give less weight to battles
	LoseOverWinRatio float64

	// WinThreshold represents the threshold upon which we consider we will surely win (for instance P > 0.8 => P = 1), 1-winThreshold represents the loseThreshold
	WinThreshold float64
}

func (hp *HeuristicParameters) String() string {
	return fmt.Sprintf(
		"c%3.2f_b%3.2f_nb%3.2f_cs%4.3f_ws%3.2e_lowr%3.2f_wt%3.2f",
		hp.Counts, hp.Battles, hp.NeutralBattles, hp.CumScore, hp.WinScore, hp.LoseOverWinRatio, hp.WinThreshold,
	)
}

// NewDefaultHeuristicParameters creates heuristic parameters
func NewDefaultHeuristicParameters() HeuristicParameters {
	return HeuristicParameters{
		Counts:           1,
		Battles:          0.5,
		NeutralBattles:   0.5,
		CumScore:         0.001,
		WinScore:         1e10,
		LoseOverWinRatio: 1.2,
		WinThreshold:     0.8,
	}
}

// Heuristic represents a heuristic
type Heuristic struct {
	HeuristicParameters
}

func (h *Heuristic) String() string {
	return h.HeuristicParameters.String()
}

// NewHeuristic
func NewHeuristic(params HeuristicParameters) Heuristic {
	return Heuristic{params}
}

// generateCoups generates coups for a given state and a given race
// While a player _can_ make multiple moves within a coup, for now this function only
// returns individual moves.
func generateCoups(s model.State, race model.Race) []model.Coup {
	all := []model.Coup{}

	for coord, cell := range s.Grid {
		if cell.Race != race {
			continue
		}

		moves := generateMovesFromCell(s, coord, cell)
		max := len(all)
		for _, move := range moves {
			// Add the move alone
			all = append(all, model.Coup{move})

			// Add the move to all the previous coups
			for _, coup := range all[:max] {
				// We have to make a copy here otherwise we will reuse the same array which will cause issues
				newCoup := make(model.Coup, len(coup)+1)
				copy(newCoup, coup)
				newCoup[len(coup)] = move

				all = append(all, newCoup)
			}
		}
	}
	return all
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

func generateMovesFromCell(s model.State, source model.Coordinates, cell model.Cell) []model.Move {
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

		moves = append(moves, model.Move{Start: source, N: cell.Count, End: target})

		// TODO: splits
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
func (h *Heuristic) scoreState(s model.State) float64 {

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
	cumScore := s.CumulativeScore

	// Win and lose cases
	if counts.ally == 0 {
		return -(h.WinScore * h.LoseOverWinRatio) + cumScore
	} else if counts.enemy == 0 {
		return h.WinScore + cumScore
	}

	for _, heuristic := range []struct {
		coeff  float64
		scores scoreCounter
	}{
		{h.Counts, counts},
		{h.Battles, battleCounts},
		{h.NeutralBattles, neutralBattleCounts},
	} {
		score := heuristic.scores.ally - heuristic.scores.enemy
		total += score * heuristic.coeff
	}

	return total + (cumScore * h.CumScore)
}
