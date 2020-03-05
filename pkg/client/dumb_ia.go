package client

import (
	"github.com/langorou/langorou/pkg/client/model"
	"log"
	"math/rand"
	"time"
)

type DumbIA struct{}

var _ IA = &DumbIA{}

func NewDumbIA() *DumbIA {
	return &DumbIA{}
}

func (dia *DumbIA) Play(state model.State) model.Coup {
	// Simulate computation
	time.Sleep(time.Second)

	coups := generateCoups(state, model.Ally)
	log.Printf("coups: %v", coups)

	if len(coups) == 0 {
		return nil
	}

	return coups[rand.Intn(len(coups))]
}
