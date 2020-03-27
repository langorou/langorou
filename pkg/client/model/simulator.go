package model

import (
	"sort"
)

type PotentialState struct {
	*State
	P float64
}

// ApplyCoup computes the possibles states after applying a Coup (a list of moves)
func (s *State) ApplyCoup(race Race, coup Coup, winThreshold float64) []PotentialState {
	// Start with the current state with probability 1
	states := []PotentialState{{State: s.Copy(true), P: 1}}

	// Sort moves by target cell
	sort.Sort(coup)

	lastEndCoordinates := coup[0].End
	var count uint8
	for _, move := range coup {
		// Aggregate by destination cell
		// This is to take into account the case where 2 groups of race A are going on cell of race B
		// We don't want two battles to occur, but only one, and the # of A should be the sum of the 2 groups

		// Move the start populations on each possible states
		for _, state := range states {
			state.DecreaseCell(move.Start, race, move.N)
		}

		// If the target cell is no more the same, stop aggregating and compute the possible states
		if move.End != lastEndCoordinates {
			states = applyMoveOnPossibleStates(states, race, lastEndCoordinates, count, winThreshold)
			count = 0
		}

		lastEndCoordinates = move.End
		count += move.N
	}

	// Apply the remaining moves
	states = applyMoveOnPossibleStates(states, race, lastEndCoordinates, count, winThreshold)

	return states
}

// applyMoveOnPossibleStates is used by applyCoup to iteratively compute the list of possible states
// that can be reached from a state and a list of moves (a coup)
func applyMoveOnPossibleStates(states []PotentialState, race Race, target Coordinates, count uint8, winThreshold float64) []PotentialState {
	// We will have at lest len(states)
	result := make([]PotentialState, 0, len(states))

	for _, state := range states {
		outcomes := applyMove(state.State, race, target, count, winThreshold)

		for _, outcome := range outcomes {
			// Take into account the probability of the previous states
			outcome.P *= state.P
			result = append(result, outcome)
		}
	}

	return result
}

// applyMove computes the possible next states from a given state and ONLY ONE move
// XXX: WARNING it re-uses the given state, so it will become stale after
func applyMove(s *State, race Race, target Coordinates, count uint8, winThreshold float64) []PotentialState {

	endCell := s.Grid[target]

	if endCell.IsEmpty() || race == endCell.Race {
		// nobody on there, or same race as ours, no battle and we can just increase the count

		// Update the cells
		s.SetCell(target, race, endCell.Count+count)
		return []PotentialState{{State: s, P: 1}}
	}

	// Fight with the enemy or neutral
	// We use a float here for later computation
	var isNeutral float64 = 0
	if endCell.Race == Neutral {
		isNeutral = 1
	}

	P := WinProbability(count, endCell.Count, isNeutral == 1)

	if P >= winThreshold {
		// Consider it a win situation given the probability
		endCount := uint8(P*float64(count) + isNeutral*float64(endCell.Count)*P)
		s.SetCell(target, race, endCount)
		return []PotentialState{{State: s, P: 1}}
	} else if P < 1-winThreshold {
		// Consider it a lose situation given the probability
		endCount := uint8((1 - P) * float64(endCell.Count))
		s.SetCell(target, endCell.Race, endCount)
		return []PotentialState{{State: s, P: 1}}
	}

	winState := s

	winState.SetCell(
		target,
		race,
		// each ally has probability P to survive. Against neutral, we have a probability P to convert them
		uint8(P*float64(count)+isNeutral*float64(endCell.Count)*P),
	)

	lossState := s.Copy(false)
	lossState.SetCell(
		target,
		endCell.Race,
		// each enemy has probability 1-P to survive
		uint8((1-P)*float64(endCell.Count)),
	)

	return []PotentialState{
		{State: winState, P: P},
		{State: lossState, P: 1 - P},
	}
}

// Adapted from github.com/Succo/twilight, but we should use float since we evaluate probability of winning.

// WinProbability of winning for the attaquant 1 with an effectif E1, agains E2
// E2 might be Neutral
func WinProbability(E1, E2 uint8, E2isNeutral bool) float64 {
	// True by property
	if (E2isNeutral && E1 >= E2) || (!E2isNeutral && float64(E1) >= 1.5*float64(E2)) {
		return 1
	}

	if E1 == E2 {
		return 0.5
	}

	if E1 < E2 {
		return float64(E1) / (2 * float64(E2))
	} else {
		return (float64(E1) / float64(E2)) - 0.5
	}
}
