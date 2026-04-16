package main

import (
	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/game/entities"
)

// add state to map and check if it change from last state
func shouldInclude[T comparable, I ~uint32](
	mode game.SnapshotMode,
	lastMap map[I]T,
	id I,
	current T,
) bool {
	last, exist := lastMap[id]
	if mode == game.SnapshotFull || !exist || last != current {
		lastMap[id] = current
		return true
	}
	lastMap[id] = current
	return false
}

func (w *World) broadcastSnapshot(snapshot []byte) {
	for id := range w.players {
		err := w.network.SendSnapshot(id, snapshot)
		if err != nil {
			// TODO: think about send snapshot error
			panic(err)
		}
	}
}

func (w *World) sendSnapshot(mode game.SnapshotMode) {
	if len(w.players) < 1 {
		return
	}

	w.currentSnapshot.Tick = w.tick
	w.currentSnapshot.Full = uint8(mode)

	w.currentSnapshot.Players = w.currentSnapshot.Players[:0]
	w.currentSnapshot.Projectiles = w.currentSnapshot.Projectiles[:0]
	w.currentSnapshot.Asteroids = w.currentSnapshot.Asteroids[:0]

	// ** Players **
	playerState := game.PlayerState{}
	for _, p := range w.players {
		playerState.ID = p.ID
		playerState.X = int32(p.Position.X * 100)
		playerState.Y = int32(p.Position.Y * 100)
		playerState.VX = int16(p.Velocity.X * 100)
		playerState.VY = int16(p.Velocity.Y * 100)
		playerState.Angle = game.RadToUint16(p.Rotation)

		if shouldInclude(mode, w.lastPlayerSnapshots, entities.PlayerID(p.Id()), playerState) {
			w.currentSnapshot.Players = append(w.currentSnapshot.Players, playerState)
		}
	}

	// ** Bullets **
	projectileState := game.ProjectileState{}
	for _, b := range w.bullets {
		projectileState.ID = b.ID
		projectileState.X = int32(b.Position.X * 100)
		projectileState.Y = int32(b.Position.Y * 100)
		projectileState.VX = int16(b.Velocity.X * 100)
		projectileState.VY = int16(b.Velocity.Y * 100)

		if shouldInclude(mode, w.lastBulletSnapshots, b.Id(), projectileState) {
			w.currentSnapshot.Projectiles = append(w.currentSnapshot.Projectiles, projectileState)
		}
	}

	// ** Asteroid **
	for _, a := range w.asteroids {
		projectileState.ID = a.ID
		projectileState.X = int32(a.Position.X * 100)
		projectileState.Y = int32(a.Position.Y * 100)
		projectileState.VX = int16(a.Velocity.X * 100)
		projectileState.VY = int16(a.Velocity.Y * 100)

		if shouldInclude(mode, w.lastAsteroidSnapshots, a.Id(), projectileState) {
			w.currentSnapshot.Asteroids = append(w.currentSnapshot.Asteroids, projectileState)
		}
	}

	n, err := w.currentSnapshot.Encode(w.buf)
	if err != nil {
		panic(err)
	}
	if len(w.currentSnapshot.Players) > 0 {
		// fmt.Printf("SENT snapshot: %+v\n", w.currentSnapshot.Players)
	}
	w.broadcastSnapshot(w.buf[:n])
}
