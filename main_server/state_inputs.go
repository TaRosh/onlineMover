package main

import (
	"github.com/TaRosh/online_mover/game"
)

func stateInputs(w *World) stateFn {
	for {
		select {
		case input := <-w.inputsQueue:
			p := w.players[input.ID]
			game.ApplyInput(p, input)
			if input.Buttons&game.InputShoot != 0 {
				if w.tick-p.LastFireTick > 10 {
					w.spawnBullet(p)
					p.LastFireTick = w.tick
				}
			}
			w.currentSnapshot.LastInputTick = input.Tick
		default:
			return stateSpawn
		}
	}
}
