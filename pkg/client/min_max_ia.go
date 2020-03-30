package client

import (
	"fmt"
	"time"

	"github.com/langorou/langorou/pkg/client/model"
)

type MinMaxIA struct {
	timeout   time.Duration
	heuristic Heuristic
}

var _ IA = &MinMaxIA{}

func NewMinMaxIA(timeout time.Duration) *MinMaxIA {
	return &MinMaxIA{
		timeout:   timeout,
		heuristic: NewHeuristic(NewDefaultHeuristicParameters()),
	}
}

func (m *MinMaxIA) Play(state *model.State) model.Coup {
	return m.heuristic.findBestCoupWithTimeout(state.Copy(false), m.timeout)
}

func (m *MinMaxIA) Name() string {
	return fmt.Sprintf("min_max_%d_%s", m.timeout, m.heuristic.String())
}
