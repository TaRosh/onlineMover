package main

import (
	"github.com/TaRosh/online_mover/game"
)

func (g *Game) updateProjectiles() {
	for _, b := range g.projectiles {
		b.Update()
	}
}

type stateFn func(g *Game) stateFn

// Update loop cut on the state
// 1. Catch packets
// 2. Reconcile + Interpolation
// 3. Get inputs
// 4. Apply inputs
// 5. Update objects in world
// 6. Input history + send inputs
// update is a one tick
func (g *Game) Update() error {
	for state := stateNetwork; state != nil; {
		state = state(g)
	}

	// Apply inputs
	// for _, p := range g.players {
	game.ApplyInput(&g.localPlayer.Player, g.input)

	// }
	// Update player with new inputs
	g.localPlayer.Update()
	for _, b := range g.projectiles {
		b.Update()
	}
	// Add to input history
	g.inputsHistory = append(g.inputsHistory, g.input)

	n, err := g.input.Encode(g.buf)
	if err != nil {
		return nil
	}

	// Send inputs

	err = g.Network.SendInput(g.buf[:n])
	if err != nil {
		panic(err)
	}
	g.Tick += 1

	return nil
}
