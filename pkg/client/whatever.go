package client

// Whatever implements the Client interface using a TCP server
type Whatever struct {
}

// NewWhatever creates a new TCP client
func NewWhatever() (*Whatever, error) {
	return &Whatever{}, nil
}

func (T Whatever) Nme() string {
	panic("implement me")
}

func (T Whatever) Mov() []Move {
	panic("implement me")
}

func (T Whatever) Set(n uint8, m uint8) error {
	panic("implement me")
}

func (T Whatever) Hum(n uint8, coords []Coordinates) error {
	panic("implement me")
}

func (T Whatever) Hme(x uint8, y uint8) error {
	panic("implement me")
}

func (T Whatever) Upd(n uint8, changes []Changes) error {
	panic("implement me")
}

func (T Whatever) Map(n uint8, changes []Changes) error {
	panic("implement me")
}

func (T Whatever) End() error {
	panic("implement me")
}

var _ Client = &Whatever{}
