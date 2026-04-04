package main

import (
	"github.com/TaRosh/online_mover/game"
	"github.com/quasilyte/gmath"
)

// find closest prev and next snapshot for given tick
func (g *Game) findNearestSnapshots(renderTick uint32) (*game.Snapshot, *game.Snapshot) {
	var s1, s2 *game.Snapshot
	for i := 0; i < len(g.snapshotBuffer)-1; i++ {
		if g.snapshotBuffer[i].Tick <= renderTick &&
			g.snapshotBuffer[i+1].Tick >= renderTick {
			s1 = g.snapshotBuffer[i]
			s2 = g.snapshotBuffer[i+1]
			break
		}
	}
	return s1, s2
}

// interpolate player position between two snapshots
func (g *Game) playerInterpolation(s1 *game.Snapshot, s2 *game.Snapshot, alpha float64, player *Player) gmath.Vec {
	if s1 == nil || s2 == nil {
		return player.Position
	}
	var posFrom, posTo gmath.Vec
	var foundS1, foundS2 bool
	for i := range s1.Players {
		if s1.Players[i].ID == player.ID {
			statePosToPlayerPos(&posFrom, s1.Players[i])
			foundS1 = true
			break
		}
	}
	for i := range s2.Players {
		if s2.Players[i].ID == player.ID {
			statePosToPlayerPos(&posTo, s2.Players[i])
			foundS2 = true
			break
		}
	}
	if !foundS1 || !foundS2 {
		return player.Position
	}
	newPos := gmath.Vec{
		X: gmath.Lerp(posFrom.X, posTo.X, alpha),
		Y: gmath.Lerp(posFrom.Y, posTo.Y, alpha),
	}
	return newPos
}
