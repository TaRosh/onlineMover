package main

import (
	"log"
	"time"

	"github.com/TaRosh/online_mover/udp"
)

const tickTime = time.Second / 20

func main() {
	ticker := time.Tick(tickTime)
	s, err := udp.NewServer("9000")
	if err != nil {
		log.Fatal(err)
	}

	incomingPackets := make(chan udp.Packet, 5)
	go s.Receive(incomingPackets)
	for range ticker {
		// run network layer
		for {
			select {
			case packet := <-incomingPackets:
				// save inputs to game
				proccessPacket(packet)
			default:
				goto END_NETWORK
			}
		}
	END_NETWORK:
		// update world on network inputs result
		// updateWorld()
		// send snapshot on network layer
		// sendSnapshot()
	}
}

func proccessPacket(packet udp.Packet) {
}
