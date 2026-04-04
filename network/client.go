package network

import (
	"github.com/TaRosh/online_mover/udp/host"
)

type client struct {
	transport host.ClientHost
}

func NewClient(h, port string) (*client, error) {
	c := client{}
	transport, err := host.NewClient(h, port)
	if err != nil {
		return nil, err
	}
	c.transport = transport
	return &c, nil
}
