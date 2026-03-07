package main

import (
	"log"
	"time"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/udp"
)

const tickTime = time.Second / 20

type World struct {
	players     []*game.Player
	inputsQueue chan game.Input
	tick        uint32
	network     udp.NetworkServer
	buf         []byte
}

func main() {
	ticker := time.Tick(tickTime)
	world := World{}
	world.players = append(world.players, game.NewPlayer())
	world.inputsQueue = make(chan game.Input, 1024)
	world.buf = make([]byte, 2048)

	s, err := udp.NewServer("9000")
	if err != nil {
		log.Fatal(err)
	}
	world.network = s

	incomingPackets := make(chan udp.Packet, 1024)
	go s.Receive(incomingPackets)
	for range ticker {
		// run network layer
		for {
			select {
			case packet := <-incomingPackets:
				// save inputs to game
				world.proccessPacket(packet)
			default:
				goto END_NETWORK
			}
		}
	END_NETWORK:
		// update world on network inputs result
		for {
			select {
			case input := <-world.inputsQueue:
				for _, p := range world.players {
					game.ApplyInput(p, input)
				}
			default:
				goto END_APPLY_INPUT
			}
		}
	END_APPLY_INPUT:
		world.update()
		// send snapshot on network layer
		world.sendSnapshot()
		world.tick += 1
	}
}

func (w *World) sendSnapshot() {
	snapshot := game.Snapshot{
		Players: make([]game.PlayerState, len(w.players)),
	}
	snapshot.Tick = w.tick
	for i, p := range w.players {
		playerState := game.PlayerState{}
		playerState.ID = uint32(i)
		playerState.Position = p.Position
		playerState.Velocity = p.Velocity
		snapshot.Players[i] = playerState
	}
	n, err := snapshot.Encode(w.buf)
	if err != nil {
		panic(err)
	}
	err = w.network.SendSnapshot(w.buf[:n])
	if err != nil {
		panic(err)
	}
}

func (w *World) update() {
	for _, p := range w.players {
		p.Update()
	}
}

func (w *World) proccessPacket(packet udp.Packet) {
	if packet.Type == udp.InputPacket {
		var input game.Input
		err := input.Decode(packet.Data)
		if err != nil {
			panic(err)
		}
		w.inputsQueue <- input
	}
}
