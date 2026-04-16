package game

func stateReconcile(g *Game) stateFn {
	if g.lastSnapshotForReconcile == nil {
		return stateInterpolation(g)
	}
	players := g.lastSnapshotForReconcile.Players
	// found our local player in snapshot
	for _, player := range players {
		if player.ID == g.localPlayer.ID {
			g.ReapplyPositionFromSnapshot(player)
			break
		}
	}
	g.lastSnapshotForReconcile = nil
	return stateInterpolation(g)
}

/*
TODO: do i need this?
	// get projectiles from snapshot
	for _, bullet := range g.lastSnapshotForReconcile.Projectiles {

		b, exist := g.projectiles[bullet.ID]
		if !exist {
			b = NewBullet(bullet.ID, color.RGBA{0, 0, 0xff, 0xff})
			g.projectiles[b.ID] = b
		}
		stateToBullet(b, bullet)
*/
