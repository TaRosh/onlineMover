package network

import (
	"github.com/TaRosh/online_mover/game"
	p "github.com/TaRosh/online_mover/udp/packet"
)

func (s *server) processPacket(packet *p.Packet, inputHere chan<- game.Input, events chan<- game.Event) {
	// decide which type
	// and what to do next
	// case 1 request for id
	// case 2 player input
	// default drop
	switch packet.Type {

	case p.Invalid:
		return

	// event for game player connect chan
	case p.Connect:
		connID := packet.ConnectionID
		exist := s.isPlayerExist(connID)
		if exist {
			// TODO: what do if old player send connect
			// for now drop
			return
		}
		id := s.newConnection(connID)
		events <- game.EventInitConnection{
			ID: id,
		}
	case p.Input:

		// for inputs -> game -> player inputs chan
		var input game.Input
		err := input.Decode(packet.Data)
		if err != nil {
			// corruption data in packet
			// just skip
			return
		}
		inputHere <- input

		// TODO: think about disconect user

	}
}
