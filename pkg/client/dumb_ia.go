package client

import (
	"log"
	"math/rand"
	"time"

	"github.com/langorou/langorou/pkg/client/model"
)

type DumbIA struct{}

var _ IA = &DumbIA{}

func NewDumbIA() *DumbIA {
	return &DumbIA{}
}

func (dia *DumbIA) Play(state *model.State) model.Coup {
	// Simulate computation
	time.Sleep(time.Second)

	return randomMove(state)
}

func (dia *DumbIA) Name() string {
	return "dumb_ia"
}

func randomMove(state *model.State) model.Coup {
	coups := generateCoups(state, model.Ally)
	log.Printf("coups: %v", coups)

	if len(coups) == 0 {
		return nil
	}

	return coups[rand.Intn(len(coups))]
}