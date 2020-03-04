package client

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const testDepth = 5

func TestMinMax(t *testing.T) {

	t.Run("case1", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10N | XXX | 06N
		// XXX | 08A | XXX
		// 12E | XXX | XXX

		startState := state{
			grid: map[Coordinates]cell{
				{X: 0, Y: 0}: {
					race:  Neutral,
					count: 10,
				},
				{X: 2, Y: 0}: {
					race:  Neutral,
					count: 6,
				},
				{X: 0, Y: 2}: {
					race:  Enemy,
					count: 12,
				},
				{X: 1, Y: 1}: {
					race:  Ally,
					count: 8,
				},
			},
			height: 3,
			width:  3,
		}

		coup, _ := minimax(startState, Ally, testDepth)
		assert.Equal(t, Coup{Move{
			Start: Coordinates{X: 1, Y: 1},
			N:     8,
			End:   Coordinates{X: 2, Y: 0},
		}}, coup)
	})

	t.Run("case2", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 08A | 06N | XXX
		// XXX | 10N | XXX
		// XXX | XXX | XXX
		// Far away: 8Enemy

		// Here our only chance to win is to try to steal the group of 10 neutrals then go for the group of 6 neutrals
		startState := state{
			grid: map[Coordinates]cell{
				{X: 0, Y: 0}: {
					race:  Ally,
					count: 8,
				},
				{X: 1, Y: 0}: {
					race:  Neutral,
					count: 6,
				},
				{X: 1, Y: 1}: {
					race:  Neutral,
					count: 10,
				},
				{X: 8, Y: 8}: {
					race:  Enemy,
					count: 8,
				},
			},
			height: 10,
			width:  10,
		}

		coup, _ := minimax(startState, Ally, testDepth)
		assert.Equal(t, Coup{Move{
			Start: Coordinates{X: 0, Y: 0},
			N:     8,
			End:   Coordinates{X: 1, Y: 0},
		}}, coup)
	})

	t.Run("case3", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// XXX | XXX | XXX | XXX
		// XXX | 68A | XXX | XXX
		// XXX | XXX | 07N | XXX
		// XXX | XXX | XXX | XXX
		// Far away: 75Enemy

		startState := state{
			grid: map[Coordinates]cell{
				{X: 1, Y: 1}: {
					race:  Ally,
					count: 68,
				},
				{X: 2, Y: 2}: {
					race:  Neutral,
					count: 7,
				},
				{X: 7, Y: 4}: {
					race:  Enemy,
					count: 75,
				},
			},
			height: 10,
			width:  10,
		}

		log.Printf("start: %+v", startState)
		coup, score := minimax(startState, Ally, testDepth)
		log.Printf("best %+v, %f", coup, score)
		assert.Equal(t, Coup{Move{
			Start: Coordinates{X: 1, Y: 1},
			N:     68,
			End:   Coordinates{X: 2, Y: 2},
		}}, coup)
	})
}

func TestSimulationAllyNeutral(t *testing.T) {
	startState := state{
		grid: map[Coordinates]cell{
			{X: 0, Y: 0}: {
				race:  Neutral,
				count: 10,
			},
			{X: 1, Y: 1}: {
				race:  Ally,
				count: 20,
			},
		},
		height: 2,
		width:  2,
	}

	t.Run("sure win", func(t *testing.T) {
		s := startState.deepCopy()

		coup := Coup{
			{Start: Coordinates{X: 1, Y: 1}, End: Coordinates{0, 0}, N: 11},
		}

		potentialStates := applyCoup(s, Ally, coup)

		assert.Len(t, potentialStates, 1)
		assert.Equal(t, potentialState{
			probability: 1,
			s: state{
				grid: map[Coordinates]cell{
					{X: 0, Y: 0}: {
						race:  Ally,
						count: 21,
					},
					{X: 1, Y: 1}: {
						race:  Ally,
						count: 9,
					},
				},
				height: 2,
				width:  2,
			},
		}, potentialStates[0])
	})

	t.Run("unsure win", func(t *testing.T) {
		s := startState.deepCopy()

		coup := Coup{
			{Start: Coordinates{X: 1, Y: 1}, End: Coordinates{0, 0}, N: 8},
		}

		potentialStates := applyCoup(s, Ally, coup)

		assert.Len(t, potentialStates, 2)
		assert.Equal(t, potentialState{
			probability: 0.4,
			s: state{
				grid: map[Coordinates]cell{
					{X: 0, Y: 0}: {
						race:  Ally,
						count: 7,
					},
					{X: 1, Y: 1}: {
						race:  Ally,
						count: 12,
					},
				},
				height: 2,
				width:  2,
			},
		}, potentialStates[0])

		assert.Equal(t, potentialState{
			probability: 0.6,
			s: state{
				grid: map[Coordinates]cell{
					{X: 0, Y: 0}: {
						race:  Neutral,
						count: 6,
					},
					{X: 1, Y: 1}: {
						race:  Ally,
						count: 12,
					},
				},
				height: 2,
				width:  2,
			},
		}, potentialStates[1])
	})

	t.Run("minmax decision", func(t *testing.T) {
		s := startState.deepCopy()
		s.grid[Coordinates{X: 1, Y: 0}] = cell{
			race:  Enemy,
			count: 15,
		}
		coup, _ := minimax(s, Ally, testDepth)

		assert.Equal(t, Coup{Move{
			Start: Coordinates{X: 1, Y: 1},
			N:     20,
			End:   Coordinates{X: 0, Y: 0},
		}}, coup)
	})
}
