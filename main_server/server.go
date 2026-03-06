package main

import (
	"fmt"
	"log"
	"time"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/udp"
	"github.com/quasilyte/gmath"
)

const tickTime = time.Second / 20

type Player struct {
	Pos gmath.Vec
	Vel gmath.Vec
	Acc gmath.Vec
}

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
	var input game.Input
	err := input.Decode(packet.Data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", input)
}
