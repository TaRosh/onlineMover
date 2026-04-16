package game

func stateInterpolation(g *Game) stateFn {
	renderTick := g.lastServerTick - g.interpolationDelay

	prevSnapshot, nextSnapshot := g.FindNearestSnapshots(renderTick)
	if prevSnapshot == nil || nextSnapshot == nil {
		return stateInput(g)
	}

	alpha := float64((renderTick - prevSnapshot.Tick)) / float64((nextSnapshot.Tick - prevSnapshot.Tick))
	for id, player := range g.players {
		// skip local player
		if id == g.localPlayer.ID {
			continue
		}
		g.players[player.ID].Position = g.PlayerInterpolation(prevSnapshot, nextSnapshot, alpha, g.players[player.ID])
		g.players[player.ID].Rotation = g.InterpolateAngle(prevSnapshot, nextSnapshot, alpha, g.players[player.ID])
	}
	return stateInput(g)
}
