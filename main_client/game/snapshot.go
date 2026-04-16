package game

import (
	"github.com/TaRosh/online_mover/game"
)

// find closest prev and next snapshot for given tick
func (g *Game) FindNearestSnapshots(renderTick uint32) (*game.Snapshot, *game.Snapshot) {
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
