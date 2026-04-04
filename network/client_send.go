package network

import "github.com/TaRosh/online_mover/udp/packet"

func (c *client) SendPlayerConnectionRequest() error {
	err := c.transport.Send(packet.Connect, nil)
	return err
}

func (c *client) SendInput(data []byte) error {
	err := c.transport.Send(packet.Input, data)
	return err
}
