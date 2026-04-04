package host

import (
	"github.com/TaRosh/online_mover/udp/packet"
)

type ServerHost interface {
	Host
	Send(connID uint32, t packet.Type, data []byte) error
	CheckTimeouts(timeoutConnectionsHere chan<- uint32)
	DeleteConn(id uint32)
}
type ClientHost interface {
	Host
	Send(t packet.Type, data []byte) error
}

type Host interface {
	Receive(sendPacketHere chan<- packet.Packet) error
	// GetAddr() *net.UDPAddr
	Close()
}
