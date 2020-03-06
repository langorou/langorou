package client

import (
	"fmt"
	"github.com/langorou/langorou/pkg/client/model"
	"sort"
)

type potentialState struct {
	s           model.State
	probability float64
}

// applyCoup computes the possibles states after applying a Coup (a list of moves)
func applyCoup(origState model.State, race model.Race, coup model.Coup) []potentialState {
	// TODO improve this function, it's not really efficient, some moves are

	// Start with the current state with probability 1
	states := []potentialState{{s: origState.Copy(true), probability: 1}}

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
			state.s.DiffDecreaseCell(move.Start, move.N)
		}

		// If the target cell is no more the same, stop aggregating and compute the possible states
		if move.End != lastEndCoordinates {

			startCell, _ := origState.GetCellWithDiff(move.Start)
			// Ensure that the move is legal
			if startCell.Race != race {
				panic(fmt.Sprintf("Race: %+v, tried to move race: %+v, illegal move", race, startCell.Race))
			}

			states = applyMoveOnPossibleStates(states, race, lastEndCoordinates, count)
			count = 0
		}

		lastEndCoordinates = move.End
		count += move.N
	}

	// Apply the remaining moves
	states = applyMoveOnPossibleStates(states, race, lastEndCoordinates, count)

	return states
}

// applyMoveOnPossibleStates is used by applyCoup to iteratively compute the list of possible states
// that can be reached from a state and a list of moves (a coup)
func applyMoveOnPossibleStates(states []potentialState, race model.Race, target model.Coordinates, count uint8) []potentialState {
	// We will have at lest len(states)
	result := make([]potentialState, 0, len(states))

	for _, state := range states {
		outcomes := applyMove(state.s, race, target, count)

		for _, outcome := range outcomes {
			// Take into account the probability of the previous states
			outcome.probability *= state.probability
			result = append(result, outcome)
		}
	}

	return result
}

// applyMove computes the possible next states from a given state and ONLY ONE move
// XXX: WARNING it re-uses the given state, so it will become stale after
func applyMove(s model.State, race model.Race, target model.Coordinates, count uint8) []potentialState {

	endCell, _ := s.GetCellWithDiff(target)

	if endCell.IsEmpty() || race == endCell.Race {
		// nobody on there, or same race as ours, no battle and we can just increase the count

		// Update the cells
		s.DiffSetCell(target, race, endCell.Count+count)
		return []potentialState{{s: s, probability: 1}}
	}

	// Fight with the enemy or neutral
	// We use a float here for later computation
	var isNeutral uint8 = 0
	if endCell.Race == model.Neutral {
		isNeutral = 1
	}

	P := winProbability(count, endCell.Count, isNeutral == 1)

	// TODO: maybe we should consider, probability > threshold as 1 as well (for instance threshold = 0.9) to lower # of computations
	if P == 1 {
		// We surely win, same as "nobody there"
		s.DiffSetCell(target, race, count+(isNeutral*endCell.Count)) // if we totally win against Neutral, we convert all of them
		return []potentialState{{s: s, probability: 1}}
	}

	winState := s

	winState.DiffSetCell(
		target,
		race,
		// each ally has probability P to survive. Against neutral, we have a probability P to convert them
		uint8(P*float64(count)+float64(isNeutral*endCell.Count)*P),
	)

	lossState := s.Copy(false)
	lossState.DiffSetCell(
		target,
		endCell.Race,
		// each enemy has probability 1-P to survive
		uint8((1-P)*float64(endCell.Count)),
	)

	return []potentialState{
		{s: winState, probability: P},
		{s: lossState, probability: 1 - P},
	}
}

// Adapted from github.com/Succo/twilight, but we should use float since we evaluate probability of winning.

// winProbability of winning for the attaquant 1 with an effectif E1, agains E2
// E2 might be Neutral
func winProbability(E1, E2 uint8, E2isNeutral bool) float64 {
	if E1 == E2 {
		return 0.5
	}

	// True by property
	if (E2isNeutral && E1 > E2) || (!E2isNeutral && float64(E1) >= 1.5*float64(E2)) {
		return 1
	}

	if E1 < E2 {
		return float64(E1) / (2 * float64(E2))
	} else {
		return (float64(E1) / float64(E2)) - 0.5
	}
}
