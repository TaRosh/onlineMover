package main

import (
	"fmt"
	"time"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/network"
)

const tickTime = time.Second / 20

type World struct {
	players     map[game.PlayerID]*game.Player
	inputsQueue chan game.Input
	// for now we have one event: receive new player id
	incomingEvents chan game.PlayerID
	snapshot       *game.Snapshot
	tick           uint32
	network        network.NetworkServer
	buf            []byte
}

func (w *World) Init() {
	w.players = make(map[game.PlayerID]*game.Player)
	w.inputsQueue = make(chan game.Input, 1024)
	w.incomingEvents = make(chan game.PlayerID, 1024)
	w.snapshot = new(game.Snapshot)
	w.tick = 0
	network, err := network.NewServer("localhost", "9000")
	if err != nil {
		// can't create udp server
		panic(err)
	}
	w.network = network
	w.buf = make([]byte, 2048)
}

func (w *World) processEvent(id game.PlayerID) {
	player := game.NewPlayer(id)
	w.players[player.ID] = player
}

func main() {
	ticker := time.Tick(tickTime)
	world := World{}
	world.Init()
	// world.players = append(world.players, game.NewPlayer())
	// now we don't know amount of players
	// world.snapshot.Players = make([]game.PlayerState, len(world.players))
	go world.network.Receive(world.inputsQueue, world.incomingEvents)
	fmt.Println("Server ready")
	for range ticker {
		// connect new players first
		for {
			select {
			case newPlayerID := <-world.incomingEvents:
				// save inputs to game
				world.processEvent(newPlayerID)
			default:
				goto END_EVENTS
			}
		}
	END_EVENTS:
		// update world on network inputs result
		for {
			select {
			case input := <-world.inputsQueue:
				for _, p := range world.players {
					if p.ID == input.ID {
						game.ApplyInput(p, input)
					}
					// TODO: think about this can be unsync in ticks
					// from clients?
					world.snapshot.LastInputTick = input.Tick
				}
			default:
				goto END_APPLY_INPUT
			}
		}
	END_APPLY_INPUT:
		world.update()
		// send snapshot on network layer
		// time.Sleep(time.Second)
		world.sendSnapshot()
		world.tick += 1
	}
}

func (w *World) sendSnapshot() {
	if len(w.players) < 1 {
		return
	}
	w.snapshot.Tick = w.tick
	w.snapshot.Players = nil
	playerState := game.PlayerState{}
	for _, p := range w.players {
		playerState.ID = p.ID
		playerState.Position = p.Position
		playerState.Velocity = p.Velocity
		w.snapshot.Players = append(w.snapshot.Players, playerState)
	}
	n, err := w.snapshot.Encode(w.buf)
	if err != nil {
		panic(err)
	}
	fmt.Println("SENT snapshot:", w.snapshot)
	w.broadcastSnapshot(w.buf[:n])
	// err = w.network.SendSnapshot(w.buf[:n])
	// if err != nil {
	// 	panic(err)
	// }
}

func (w *World) broadcastSnapshot(snapshot []byte) {
	fmt.Println(w.players)
	for id := range w.players {
		err := w.network.SendSnapshot(id, snapshot)
		if err != nil {
			// TODO: think about send snapshot error
			panic(err)
		}
	}
}

func (w *World) update() {
	for _, p := range w.players {
		p.Update()
	}
}
