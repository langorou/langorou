package client

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

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

	// MaxGroups indicates the maximum groups of ally units we want
	MaxGroups uint8

	// Groups is used to penalize/reward the fact of having a lot of scattered units
	// We want it to be negative since we want to penalize the fact of having a lot of scattered units
	Groups float64
}

func (hp *HeuristicParameters) String() string {
	return fmt.Sprintf("%+v", *hp)
}

// NewDefaultHeuristicParameters creates defaultns heuristic parameters
func NewDefaultHeuristicParameters() HeuristicParameters {
	return HeuristicParameters{
		Counts:           1,
		Battles:          0.0,
		NeutralBattles:   0.0,
		CumScore:         0.0001,
		WinScore:         1e10,
		LoseOverWinRatio: 1,
		WinThreshold:     1.,
		MaxGroups:        2,
		Groups:           -0,
	}
}

var coupPool = sync.Pool{
	New: func() interface{} {
		return model.Coup{}
	},
}

var coupsPool = sync.Pool{
	New: func() interface{} {
		return []model.Coup{}
	},
}

func putCoup(coup model.Coup) {
	coupPool.Put(coup)
}

func putCoups(coups []model.Coup) {
	coupsPool.Put(coups)
}

func getCoup() model.Coup {
	coup := coupPool.Get()

	return coup.(model.Coup)[:0]
}

func getCoups() []model.Coup {
	coups := coupsPool.Get()

	return coups.([]model.Coup)[:0]
}

// Heuristic represents a heuristic
type Heuristic struct {
	HeuristicParameters
}

func (h *Heuristic) String() string {
	return h.HeuristicParameters.String()
}

// NewHeuristic creates a new heuristic given parameters
func NewHeuristic(params HeuristicParameters) Heuristic {
	return Heuristic{params}
}

// randomMove gives a random move among the possible moves for the race Ally
func (h *Heuristic) randomMove(state *model.State) model.Coup {
	coups := h.generateCoups(state, model.Ally)

	if len(coups) == 0 {
		return nil
	}

	return coups[rand.Intn(len(coups))]
}

// generateCoups generates coups for a given state and a given race
// While a player _can_ make multiple moves within a coup, for now this function only
// returns individual moves.
// It computes a product of all the possible moves for each group of our race (including the move that consists in not moving)
func (h *Heuristic) generateCoups(s *model.State, race model.Race) []model.Coup {
	// TODO: try pre allocating here
	all := getCoups()

	for coord, cell := range s.Grid {
		if cell.Race != race {
			continue
		}

		splitThreshold := uint8(0)
		allowSplit := (race == model.Ally && s.AlliesGroups < h.MaxGroups) || (race == model.Enemy && s.EnemiesGroups < h.MaxGroups)
		if allowSplit {
			splitThreshold = 2 * s.SmallestNeutralGroup
		}

		moves := generateMovesFromCell(s.Width, s.Height, coord, cell, splitThreshold)
		max := len(all)
		for _, move := range moves {
			// Add the move alone
			all = append(all, model.Coup{move})

			// Add the move to all the previous coups
			for _, coup := range all[:max] {
				// We have to make a copy here otherwise we will reuse the same array which will cause issues
				newCoup := getCoup()
				copy(newCoup, coup)
				newCoup = append(newCoup, move)

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

func transform(width, height uint8, c model.Coordinates, t transformation) (res model.Coordinates, ok bool) {
	if width == 0 {
		return c, false
	}

	xRes := uint8(int8(c.X) + t.X)
	if (xRes < 0) || (xRes >= width) {
		return c, false
	}

	yRes := uint8(int8(c.Y) + t.Y)
	if (yRes < 0) || (yRes >= height) {
		return c, false
	}

	return model.Coordinates{X: xRes, Y: yRes}, true
}

func generateMovesFromCell(width, height uint8, source model.Coordinates, cell model.Cell, splitThreshold uint8) []model.Move {
	// Simplification: as long as we are doing one move per turn, we are better
	// off moving all units and not a subset.
	// This is not true anymore if we do multiple moves per turn, as keeping
	// some units in the source cell allows us to do another attack (on another
	// cell) in the same turn.

	// 8 for regular moves + 8 for possible split
	moves := make([]model.Move, 0, 16)

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
		target, ok := transform(width, height, source, t)
		if !ok {
			continue
		}

		moves = append(moves, model.Move{Start: source, N: cell.Count, End: target})

		// TODO: for now we only move one subset of units, but we could move two, for instance:
		// ---
		// Currently we can't do: 24 at (0, 0) -> 12 at (1, 1) and 12 at (1, 0)
		// We can only do		  24 at (0, 0) -> 12 at (0, 0) and 12 at (1, 0)
		// ---
		// Allow to split only if we are among a threshold and we always split in 2
		if splitThreshold != 0 && cell.Count >= splitThreshold {
			moves = append(moves, model.Move{Start: source, N: cell.Count / 2, End: target})
		}
	}
	return moves
}

// scoreNeutralBattle scores the issue of a battle between a monster and a neutral group
func scoreNeutralBattle(c1, c2 model.Coordinates, cell1, cell2 model.Cell) float64 {
	proba := model.WinProbability(cell1.Count, cell2.Count, true)
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
	p1 := model.WinProbability(cell1.Count, cell2.Count, false)
	// p2 is for 2 attacks 1
	p2 := model.WinProbability(cell2.Count, cell1.Count, false)

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
func (h *Heuristic) scoreState(s *model.State) float64 {

	// different counts participating in the heuristic
	counts := scoreCounter{}
	battleCounts := scoreCounter{}
	neutralBattleCounts := scoreCounter{}

	for c1, cell1 := range s.Grid {
		if cell1.Race == model.Neutral {
			continue
		}
		counts.add(cell1.Race, float64(cell1.Count))

		// Avoid computing battles scores if the coefficients are 0
		if h.Battles == 0 && h.NeutralBattles == 0 {
			continue
		}

		// Loop to compute stats on the possible battle
		for c2, cell2 := range s.Grid {
			if c1 == c2 || cell1.Race == cell2.Race {
				continue
			}

			if cell2.Race == model.Neutral && h.NeutralBattles != 0 {
				// TODO: average here since we can count a battle multiple times, for now we just consider it as multiple opportunities, hence there is no average
				neutralBattleCounts.add(cell1.Race, scoreNeutralBattle(c1, c2, cell1, cell2))
			} else if cell2.Race == cell1.Race.Opponent() && h.Battles != 0 {
				// TODO: average here since we can count a battle multiple times, for now we just consider it as multiple opportunities, hence there is no average
				g1, g2 := scoreMonsterBattle(c1, c2, cell1, cell2)
				battleCounts.add(cell1.Race, g1)
				battleCounts.add(cell2.Race, g2)
			}
		}
	}

	total := 0.

	// TODO: try distance power alpha instead of distance power 1, caveat: computations
	cumScore := s.CumulativeScore

	// Win and lose cases
	if counts.ally == 0 {
		return -(h.WinScore * h.LoseOverWinRatio) + cumScore
	} else if counts.enemy == 0 {
		// In case of win we want to win the earliest we can, so we SUBSTRACT the cumulative score
		return h.WinScore - cumScore
	}

	groupsCounts := scoreCounter{ally: float64(s.AlliesGroups), enemy: float64(s.EnemiesGroups)}

	for _, heuristic := range []struct {
		coef   float64
		scores scoreCounter
	}{
		{h.Counts, counts},
		{h.Battles, battleCounts},
		{h.NeutralBattles, neutralBattleCounts},
		{h.Groups, groupsCounts},
	} {
		score := heuristic.scores.ally - heuristic.scores.enemy
		total += score * heuristic.coef
	}

	return total + (cumScore * h.CumScore)
}
