package host

import (
	inner "github.com/TaRosh/online_mover/udp/connection"
	p "github.com/TaRosh/online_mover/udp/packet"
)

func middlewareChangeAddr(conn *inner.Conn, packet *p.Packet) error {
	if packet.Addr == nil {
		return nil
	}
	conn.SetAddr(packet.Addr)
	return nil
}
