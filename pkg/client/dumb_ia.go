package client

import (
	"time"

	"github.com/langorou/langorou/pkg/client/model"
)

type DumbIA struct {
	h Heuristic
}

var _ IA = &DumbIA{}

func NewDumbIA() *DumbIA {
	h := NewHeuristic(NewDefaultHeuristicParameters())
	return &DumbIA{h}
}

func (dia *DumbIA) Play(state *model.State) model.Coup {
	// Simulate computation
	time.Sleep(time.Second)

	return dia.h.randomMove(state)
}

func (dia *DumbIA) Name() string {
	return "dumb_ia"
}
