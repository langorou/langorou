package client

import (
	"math"

	"github.com/langorou/langorou/pkg/client/model"
)

const (
	posInfinity = math.MaxFloat64
	negInfinity = -math.MaxFloat64
)

func (h *Heuristic) findBestCoup(state model.State, depth uint8) (model.Coup, float64) {
	return h.alphabeta(state, model.Ally, negInfinity, posInfinity, depth)
}

// alphabeta computes the best coup going at most at depth depth
func (h *Heuristic) alphabeta(state model.State, race model.Race, alpha float64, beta float64, depth uint8) (model.Coup, float64) {
	bestCoup := model.Coup{}

	// Max depth reached
	if depth <= 0 {
		return bestCoup, h.scoreState(state)
	}

	coups := generateCoups(state, race)
	// No moves found
	if len(coups) == 0 {
		return bestCoup, h.scoreState(state)
	}

	// Chose if we want to maximize (us) or minimize (enemy) our score
	value := negInfinity
	f := math.Max
	if race == model.Enemy {
		value = posInfinity
		f = math.Min
	}

	// Sort by killer moves
	model.SortCoupsByQuickScore(coups, state)

	// for each generated coup, we compute the list of potential outcomes and compute an average score
	// weighted by the probabilities of these potential outcomes
	for _, coup := range coups {
		outcomes := applyCoup(state, race, coup, h.WinThreshold)

		score := 0.
		// log.Printf("depth: %d", depth)
		for _, outcome := range outcomes {
			_, tmpScore := h.alphabeta(outcome.s, race.Opponent(), alpha, beta, depth-1)
			score += tmpScore * outcome.probability
		}

		// log.Printf("cumulative: %f, findBestCoup score: %f at depth: %d for race: %v and coup: %+v, grid: %+v, potential: %+v", state.cumScore, score, depth, race, coup, state.grid, outcomes)

		if f(value, score) == score { // score >= value if max playing or value >= score if min playing
			value = score
			bestCoup = coup
			// log.Printf("better value found %f: depth: %d, race: %v", value, depth, race)
		}

		// Check for possible cuts
		if race == model.Enemy {
			// alpha cut
			if alpha >= value {
				return bestCoup, value
			}
			beta = f(beta, value)
		} else {
			// beta cut
			if value >= beta {
				return bestCoup, value
			}
			alpha = f(alpha, value)
		}
	}

	return bestCoup, value
}
