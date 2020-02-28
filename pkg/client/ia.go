package client

type IA interface {
	Play(state state) Coup
}
