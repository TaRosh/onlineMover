package network

import (
	"log"

	p "github.com/TaRosh/online_mover/udp/packet"
	"github.com/TaRosh/online_mover/game"
)

// i think for now will use chan for each event from client
func (s *server) Receive(inputHere chan<- game.Input, events chan<- game.Event) {
	sendPacketHere := make(chan p.Packet, 1024)
	go func() {
		for {
			err := s.transport.Receive(sendPacketHere)
			if err != nil {
				log.Println(err)
			}
		}
	}()
	for packet := range sendPacketHere {
		s.processPacket(&packet, inputHere, events)
	}
}
