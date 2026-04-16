package game

import "github.com/TaRosh/online_mover/game"

func stateNetwork(g *Game) stateFn {
	const maxSnapshotsPerFrame = 5
	var processed int
	// catch packets
	for processed < maxSnapshotsPerFrame {
		select {
		case snapshot := <-g.snapshotQueue:
			// ignore old snapshot
			if snapshot.Tick < g.lastServerTick {
				continue
			}
			// make place in snapshot buffer
			if len(g.snapshotBuffer) > g.maxSnapshot {
				g.snapshotBuffer = g.snapshotBuffer[1:]
			}
			g.snapshotBuffer = append(g.snapshotBuffer, &snapshot)
			g.lastServerTick = snapshot.Tick
			g.lastSnapshotForReconcile = &snapshot

			// TODO: think abou
			if snapshot.Full == uint8(game.SnapshotFull) {
				g.HandleFullSnapshot(&snapshot)
			} else {
				g.HandleDeltaSnapshot(&snapshot)
			}
			processed += 1
		default:
			// no snapshots than finish
			processed = maxSnapshotsPerFrame
		}
	}
	return stateReconcile(g)
}
