package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/langorou/langorou/pkg/client/model"
	"github.com/stretchr/testify/assert"
)

const testDepth = 5
const testTimeout time.Duration = 1500 * time.Millisecond

func testedFindCoup(t *testing.T, state *model.State) model.Coup {
	coupDepth, _ := testHeuristic.findBestCoup(state, testDepth)
	coupTimeout := testHeuristic.findBestCoupWithTimeout(state, testTimeout)

	assert.Equal(t, coupDepth, coupTimeout)
	return coupDepth
}

var testHeuristic = NewHeuristic(NewDefaultHeuristicParameters())

func TestMinMax(t *testing.T) {

	t.Run("case1", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10N | XXX | 06N
		// XXX | 08A | XXX
		// 12E | XXX | XXX

		startState := model.NewState(3, 3)
		startState.SetCell(model.Coordinates{}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 2}, model.Neutral, 6)
		startState.SetCell(model.Coordinates{Y: 2}, model.Enemy, 12)
		startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 8)

		coup := testHeuristic.findBestCoupWithTimeout(startState, testTimeout)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 1, Y: 1},
			N:     8,
			// Or End: model.Coordinates{X: 1, Y:0 },
			End: model.Coordinates{X: 0, Y: 2},
		}}, coup)
	})

	t.Run("case1.1", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10N | XXX | 14E
		// ... | XXX | XXX
		// 12A | XXX | XXX

		startState := model.NewState(3, 3)
		startState.SetCell(model.Coordinates{}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 2}, model.Enemy, 14)
		startState.SetCell(model.Coordinates{X: 0, Y: 2}, model.Ally, 12)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 0, Y: 2},
			N:     12,
			End:   model.Coordinates{X: 0, Y: 1},
		}}, coup)
	})

	t.Run("case1.1.1", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10N | XXX | 14A
		// 12E | XXX | XXX
		// XXX | XXX | XXX

		startState := model.NewState(3, 3)
		startState.SetCell(model.Coordinates{}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 2}, model.Ally, 14)
		startState.SetCell(model.Coordinates{Y: 1}, model.Enemy, 12)

		coup := testedFindCoup(t, startState)
		assert.Len(t, coup, 1)
		assert.Equal(t, model.Coordinates{X: 2}, coup[0].Start)
		assert.EqualValues(t, 14, coup[0].N)
		// Allow both X: 1, Y:0 and X: 1, Y:1
		assert.EqualValues(t, 1, coup[0].End.X)
	})

	t.Run("case1.1.2", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10N | 14E | XXX
		// 12A | XXX | XXX
		// XXX | XXX | XXX

		startState := model.NewState(3, 3)
		startState.SetCell(model.Coordinates{}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 1}, model.Enemy, 14)
		startState.SetCell(model.Coordinates{Y: 1}, model.Ally, 12)

		// Probability  3/7 of winning if we attack the 14 group
		// Probability 15/22 of winning if we are being attacked by the 14 group when we are 22
		// So we should take the villagers
		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{Y: 1},
			N:     12,
			End:   model.Coordinates{X: 0, Y: 0},
		}}, coup)
	})

	t.Run("case1.2", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10N | XXX | 14A
		// XXX | 12E | XXX
		// XXX | XXX | XXX

		startState := model.NewState(3, 3)
		startState.SetCell(model.Coordinates{}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 2}, model.Ally, 14)
		startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Enemy, 12)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 2, Y: 0},
			N:     14,
			End:   model.Coordinates{X: 1, Y: 1},
		}}, coup)
	})

	t.Run("case1.3", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10N | XXX | 06N
		// XXX | XXX | XXX
		// 12A | 08E | XXX

		startState := model.NewState(3, 3)
		startState.SetCell(model.Coordinates{}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 2}, model.Neutral, 6)
		startState.SetCell(model.Coordinates{Y: 2}, model.Ally, 12)
		startState.SetCell(model.Coordinates{X: 1, Y: 2}, model.Enemy, 8)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 0, Y: 2},
			N:     12,
			End:   model.Coordinates{X: 1, Y: 2},
		}}, coup)
	})

	t.Run("case2", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 08A | 06N | XXX
		// XXX | 10N | XXX
		// XXX | XXX | XXX
		// Far away: 8Enemy

		// Here our only chance to win is to try to steal the group of 10 neutrals then go for the group of 6 neutrals
		startState := model.NewState(10, 10)
		startState.SetCell(model.Coordinates{}, model.Ally, 8)
		startState.SetCell(model.Coordinates{X: 1}, model.Neutral, 6)
		startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 8, Y: 8}, model.Enemy, 8)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 0, Y: 0},
			N:     8,
			End:   model.Coordinates{X: 1, Y: 0},
		}}, coup)
	})

	t.Run("case3", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// XXX | XXX | XXX | XXX
		// XXX | 68A | XXX | XXX
		// XXX | XXX | 07N | XXX
		// XXX | XXX | XXX | XXX
		// Far away: 75Enemy

		startState := model.NewState(10, 10)
		startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 68)
		startState.SetCell(model.Coordinates{X: 2, Y: 2}, model.Neutral, 7)
		startState.SetCell(model.Coordinates{X: 7, Y: 4}, model.Enemy, 75)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 1, Y: 1},
			N:     68,
			End:   model.Coordinates{X: 2, Y: 2},
		}}, coup)
	})

	t.Run("case4", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10E | XXX | XXX
		// 06N | 10A | XXX
		// XXX | XXX | XXX

		startState := model.NewState(10, 10)
		startState.SetCell(model.Coordinates{}, model.Enemy, 10)
		startState.SetCell(model.Coordinates{Y: 1}, model.Neutral, 6)
		startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 10)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 1, Y: 1},
			N:     10,
			End:   model.Coordinates{X: 0, Y: 1},
		}}, coup)
	})

	t.Run("case4.1", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10A | XXX | XXX
		// 06N | 10E | XXX
		// XXX | XXX | XXX

		startState := model.NewState(10, 10)
		startState.SetCell(model.Coordinates{}, model.Ally, 10)
		startState.SetCell(model.Coordinates{Y: 1}, model.Neutral, 6)
		startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Enemy, 10)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 0, Y: 0},
			N:     10,
			End:   model.Coordinates{X: 0, Y: 1},
		}}, coup)
	})

	t.Run("case5", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10N | XXX
		// 15E | 20A

		startState := model.NewState(2, 2)
		startState.SetCell(model.Coordinates{}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 20)
		startState.SetCell(model.Coordinates{X: 1, Y: 0}, model.Enemy, 15)

		// Probability 5/6 of winning if we attack the enemy directly
		// only 3/4 if we get the villagers but the enemy attacks us after
		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 1, Y: 1},
			N:     20,
			End:   model.Coordinates{X: 1, Y: 0},
		}}, coup)
	})

	t.Run("case6", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// XXX | XXX | 12E
		// XXX | XXX | XXX
		// 30A | XXX | XXX

		startState := model.NewState(3, 3)
		startState.SetCell(model.Coordinates{Y: 2}, model.Ally, 30)
		startState.SetCell(model.Coordinates{X: 2}, model.Enemy, 12)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 0, Y: 2},
			N:     30,
			End:   model.Coordinates{X: 1, Y: 1},
		}}, coup)
	})

	t.Run("case6.1", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// XXX | XXX | 12A
		// XXX | 30E | XXX
		// XXX | XXX | XXX

		startState := model.NewState(3, 3)
		startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Enemy, 30)
		startState.SetCell(model.Coordinates{X: 2}, model.Ally, 12)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 2, Y: 0},
			N:     12,
			End:   model.Coordinates{X: 1, Y: 1},
		}}, coup)
	})

	t.Run("case6.2", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// XXX | XXX | 12A
		// XXX | XXX | XXX
		// XXX | 30E | XXX

		startState := model.NewState(3, 3)
		startState.SetCell(model.Coordinates{X: 1, Y: 2}, model.Enemy, 30)
		startState.SetCell(model.Coordinates{X: 2}, model.Ally, 12)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 2, Y: 0},
			N:     12,
			End:   model.Coordinates{X: 1, Y: 0},
		}}, coup)
	})

	t.Run("case6.2.1", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// XXX | 12E | XXX
		// XXX | XXX | XXX
		// XXX | 30A | XXX

		startState := model.NewState(3, 3)
		startState.SetCell(model.Coordinates{X: 1, Y: 2}, model.Ally, 30)
		startState.SetCell(model.Coordinates{X: 1}, model.Enemy, 12)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 1, Y: 2},
			N:     30,
			End:   model.Coordinates{X: 1, Y: 1},
		}}, coup)
	})

	t.Run("case7", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// XXX | XXX | XXX | 12E | 02N
		// Enemy far away

		startState := model.NewState(5, 7)
		startState.SetCell(model.Coordinates{X: 3, Y: 2}, model.Ally, 12)
		startState.SetCell(model.Coordinates{X: 4, Y: 2}, model.Neutral, 2)
		startState.SetCell(model.Coordinates{Y: 5}, model.Enemy, 10)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 3, Y: 2},
			N:     12,
			End:   model.Coordinates{X: 4, Y: 2},
		}}, coup)
	})

	t.Run("case8", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10N | XXX | XXX | XXX | 10N |
		// 10A | XXX | XXX | XXX | 10A |
		// XXX | XXX | XXX | XXX | XXX |
		// 40 Enemy far far away

		startState := model.NewState(5, 20)
		startState.SetCell(model.Coordinates{X: 0, Y: 0}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 0, Y: 1}, model.Ally, 10)
		startState.SetCell(model.Coordinates{X: 4, Y: 0}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 4, Y: 1}, model.Ally, 10)
		startState.SetCell(model.Coordinates{X: 2, Y: 19}, model.Enemy, 40)

		coup := testedFindCoup(t, startState)
		assert.Equal(t, model.Coup{
			model.Move{
				Start: model.Coordinates{X: 0, Y: 1},
				N:     10,
				End:   model.Coordinates{X: 0, Y: 0},
			},
			model.Move{
				Start: model.Coordinates{X: 4, Y: 1},
				N:     10,
				End:   model.Coordinates{X: 4, Y: 0},
			},
		}, coup)
	})
}

func TestSimulationAllyNeutral(t *testing.T) {
	startState := model.NewState(2, 2)
	startState.SetCell(model.Coordinates{}, model.Neutral, 10)
	startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 20)

	t.Run("sure win", func(t *testing.T) {
		s := startState.Copy(false)

		coup := model.Coup{
			{Start: model.Coordinates{X: 1, Y: 1}, End: model.Coordinates{}, N: 11},
		}

		potentialStates := s.ApplyCoup(model.Ally, coup, 1)

		assert.Len(t, potentialStates, 1)
		expected := model.NewState(2, 2)
		expected.SetCell(model.Coordinates{}, model.Ally, 21)
		expected.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 9)
		assert.EqualValues(t, 1, potentialStates[0].P)
		assert.Equal(t, expected.Grid, potentialStates[0].Grid)
	})

	t.Run("unsure win", func(t *testing.T) {
		s := startState.Copy(false)

		coup := model.Coup{
			{Start: model.Coordinates{X: 1, Y: 1}, End: model.Coordinates{}, N: 8},
		}

		potentialStates := s.ApplyCoup(model.Ally, coup, 1)

		assert.Len(t, potentialStates, 2)
		expected1 := model.NewState(2, 2)
		expected1.SetCell(model.Coordinates{}, model.Ally, 7)
		expected1.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 12)
		assert.Equal(t, 0.4, potentialStates[0].P)
		assert.Equal(t, expected1.Grid, potentialStates[0].Grid)

		expected2 := model.NewState(2, 2)
		expected2.SetCell(model.Coordinates{}, model.Neutral, 6)
		expected2.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 12)
		assert.EqualValues(t, 0.6, potentialStates[1].P)
		assert.Equal(t, expected2.Grid, potentialStates[1].Grid)
	})
}

func BenchmarkAB(b *testing.B) {
	cplxState := model.GenerateComplicatedState()

	for _, depth := range []uint8{2, 3, 4, 5} {
		b.Run(fmt.Sprintf("depth%d_complex", depth), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				testHeuristic.findBestCoup(cplxState.Copy(false), depth)
			}
		})
	}

	b.Run("timeout2s_complex", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			testHeuristic.findBestCoupWithTimeout(cplxState.Copy(false), 2*time.Second)
		}
	})

	smplState := model.GenerateSimpleState()

	for _, depth := range []uint8{2, 3, 4, 5, 6} {
		b.Run(fmt.Sprintf("depth%d_simple", depth), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				testHeuristic.findBestCoup(smplState.Copy(false), depth)
			}
		})
	}

	b.Run("timeout2s_simple", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			testHeuristic.findBestCoupWithTimeout(smplState.Copy(false), 2*time.Second)
		}
	})
}

func BenchmarkHeuristic(b *testing.B) {
	startState := model.GenerateComplicatedState()
	for n := 0; n < b.N; n++ {
		testHeuristic.scoreState(startState)
	}
}

func BenchmarkHash(b *testing.B) {
	startState := model.GenerateComplicatedState()

	for n := 0; n < b.N; n++ {
		startState.Hash(model.Ally, testHeuristic.hashBuffer)
	}
}

func TestGenerateCoups(t *testing.T) {
	startState := model.GenerateComplicatedState()
	testHeuristic.generateCoups(startState, model.Ally)
}

func TestFindBestCoupWithTimeoutEnds(t *testing.T) {
	startState := model.GenerateSimpleState()

	for i := 0; i < 3; i++ {
		s := time.Now()
		testHeuristic.findBestCoupWithTimeout(startState, testTimeout)
		e := time.Now()

		assert.WithinDuration(t, s.Add(testTimeout), e, 50*time.Millisecond)
	}
}

func TestGenerateCoupsDifferentStartAndEndCells(t *testing.T) {
	// RULE 5 of the Game
	startState := model.NewState(2, 2)
	startState.SetCell(model.Coordinates{}, model.Ally, 1)
	startState.SetCell(model.Coordinates{X: 1}, model.Ally, 1)

	// 01A | 01A
	// XXX | XXX

	// Individual moves: 3 moves possible for each groups -> 6 Moves
	// Group moves: 4 possibles (2 groups go down, or they cross in an X movement) and 2 where they unit (bottom left or bottom right)

	coups := testHeuristic.generateCoups(startState, model.Ally)
	assert.Len(t, coups, 10)
}
