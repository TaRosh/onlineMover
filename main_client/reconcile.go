package main

import (
	"github.com/TaRosh/online_mover/game"
	"github.com/quasilyte/gmath"
)

// Set possiton for player from snapshot
// then reapply inputs from inputs history by tick identifier
// until snapshot last tick input field
func (g *Game) reapplyPossitionFromSnapshot(player game.PlayerState) {
	newPos := gmath.Vec{}
	statePosToPlayerPos(&newPos, player)
	g.localPlayer.Position = newPos
	// player.Velocity = g.lastSnapshotForReconcile.Players[0].Velocity
	inputsAfterSnapshot := g.inputsHistory[:0]
	for _, input := range g.inputsHistory {
		if input.Tick > g.lastSnapshotForReconcile.LastInputTick {
			inputsAfterSnapshot = append(inputsAfterSnapshot, input)
		}
	}
	g.inputsHistory = inputsAfterSnapshot

	for _, input := range g.inputsHistory {
		game.ApplyInput(&g.localPlayer.Player, input)
	}
	// g.debugPlayer.Position = *&player.Position
	// fmt.Println("DEBUG PLAYER POS", g.debugPlayer.Position)
}
