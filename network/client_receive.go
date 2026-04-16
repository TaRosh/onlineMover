package network

import (
	"log"

	"github.com/TaRosh/online_mover/game"
	p "github.com/TaRosh/online_mover/udp/packet"
)

func (c *client) Receive(snapshots chan<- game.Snapshot, events chan<- game.Event) {
	sendPacketHere := make(chan p.Packet, 1024)
	go func() {
		for {
			err := c.transport.Receive(sendPacketHere)
			// TODO: think about case when error from receive
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}()
	for packet := range sendPacketHere {
		// fmt.Println("RECEIVE packet", packet.String())
		c.processPacket(&packet, snapshots, events)
	}
}
