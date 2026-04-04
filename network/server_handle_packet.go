package network

import (
	"fmt"

	p "github.com/TaRosh/online_mover/udp/packet"
	"github.com/TaRosh/online_mover/game"
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
		events <- game.Event{
			Type: game.EventConnection,
			ID:   id,
		}
		playerPacketWithID := PlayerIDPacket{
			ID: uint32(id),
		}
		n, err := playerPacketWithID.Encode(s.sentBuf)
		if err != nil || n == 0 {
			// TODO: think about it
			panic(err)
		}

		fmt.Println("write to player", id)
		err = s.sentPlayerID(connID, s.sentBuf[:n])
		// TODO: think what should do?
		// maybe resend becouse player don't get his id
		// but we allready add him to our game
		// some backup plan to remove if can't sand id
		if err != nil {
			panic(err)
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
