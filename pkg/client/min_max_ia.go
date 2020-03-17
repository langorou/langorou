package client

import (
	"fmt"
	"log"

	"github.com/langorou/langorou/pkg/client/model"
)

type MinMaxIA struct {
	depth     uint8
	heuristic Heuristic
}

var _ IA = &MinMaxIA{}

func NewMinMaxIA(depth uint8) *MinMaxIA {
	return &MinMaxIA{
		depth:     depth,
		heuristic: NewHeuristic(NewDefaultHeuristicParameters()),
	}
}

func (m *MinMaxIA) Play(state *model.State) model.Coup {
	coup, score := m.heuristic.findBestCoup(state, m.depth)
	log.Printf("MinMaxIA computed a coup with score: %f", score)
	return coup
}

func (m *MinMaxIA) Name() string {
	return fmt.Sprintf("min_max_%d_%s", m.depth, m.heuristic.String())
}
