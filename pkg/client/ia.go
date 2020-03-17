package client

import "github.com/langorou/langorou/pkg/client/model"

type IA interface {
	Play(state *model.State) model.Coup
	Name() string
}
