package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/udp"
)

const tickTime = time.Second / 20

type World struct {
	players []*game.Player
}

func main() {
	ticker := time.Tick(tickTime)
	world := World{}
	world.players = append(world.players, game.NewPlayer(color.White))

	s, err := udp.NewServer("9000")
	if err != nil {
		log.Fatal(err)
	}

	incomingPackets := make(chan udp.Packet, 5)
	incomingInputs := make(chan game.Input, 5)
	go s.Receive(incomingPackets)
	for range ticker {
		// run network layer
		for {
			select {
			case packet := <-incomingPackets:
				// save inputs to game
				proccessPacket(packet, incomingInputs)
			default:
				goto END_NETWORK
			}
		}
	END_NETWORK:
		// update world on network inputs result
		game.ApplyInput(world.players[0])
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
