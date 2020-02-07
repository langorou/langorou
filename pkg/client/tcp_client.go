package client

// TCPClient implements the Client interface using a TCP server
type TCPClient struct {

}

// NewTCPClient creates a new TCP client
func NewTCPClient() (*TCPClient, error) {
	return &TCPClient{}, nil
}

func (T TCPClient) Nme() string {
	panic("implement me")
}

func (T TCPClient) Mov() []Move {
	panic("implement me")
}

func (T TCPClient) Set(n uint8, m uint8) error {
	panic("implement me")
}

func (T TCPClient) Hum(n uint8, coords []Coordinates) error {
	panic("implement me")
}

func (T TCPClient) Hme(x uint8, y uint8) error {
	panic("implement me")
}

func (T TCPClient) Upd(n uint8, changes []Changes) error {
	panic("implement me")
}

func (T TCPClient) Map(n uint8, changes []Changes) error {
	panic("implement me")
}

func (T TCPClient) End() error {
	panic("implement me")
}

var _ Client = &TCPClient{}
