package main

import (
	"fmt"
	"time"

	"github.com/TaRosh/online_mover/game"
	"github.com/joho/godotenv"
)

const tickTime = time.Second / 20

func main() {
	ticker := time.Tick(tickTime)
	world := World{}
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	world.Init()
	go world.network.Receive(world.inputsQueue, world.incomingEvents)
	fmt.Println("Server ready")
	for range ticker {

		// 1. Event
		// 2. Apply inputs
		// 3. Spawn meteors
		for state := stateEvents; state != nil; {
			state = state(&world)
		}

		// 3. Update world
		world.update()
		// TODO: add asteroid spawn

		// send snapshot on network layer
		// time.Sleep(time.Second)
		// world.sendSnapshot(SnapshotDelta)
		if world.tick-world.lastFullSnapshotTick > 2 {
			world.sendSnapshot(game.SnapshotFull)
			world.lastFullSnapshotTick = world.tick
		} else {
			world.sendSnapshot(game.SnapshotDelta)
		}
		world.tick += 1
		world.network.CheckTimeouts(world.incomingEvents)
	}
}
