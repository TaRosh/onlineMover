package game

import (
	"math"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/main_client/converter"
	"github.com/TaRosh/online_mover/main_client/player"
	"github.com/quasilyte/gmath"
)

func (g *Game) InterpolateAngle(prevSnapshot *game.Snapshot, nextSnapshot *game.Snapshot, alpha float64, player *player.Player) gmath.Rad {
	if prevSnapshot == nil || nextSnapshot == nil {
		return player.Rotation
	}
	var angleFrom, angleTo gmath.Rad
	var found1, found2 bool
	for i := range prevSnapshot.Players {
		if prevSnapshot.Players[i].ID == player.ID {
			converter.StateAngleToPlayerRotation(&angleFrom, prevSnapshot.Players[i])
			found1 = true
			break
		}
	}
	for i := range nextSnapshot.Players {
		if nextSnapshot.Players[i].ID == player.ID {
			converter.StateAngleToPlayerRotation(&angleTo, nextSnapshot.Players[i])
			found2 = true
			break
		}
	}
	if !found1 || !found2 {
		return player.Rotation
	}
	diff := math.Mod(float64(angleTo)-float64(angleFrom)+math.Pi, 2*math.Pi) - math.Pi
	return gmath.Rad(float64(angleFrom) + diff*alpha)
}

// interpolate player position between two snapshots
func (g *Game) PlayerInterpolation(s1 *game.Snapshot, s2 *game.Snapshot, alpha float64, player *player.Player) gmath.Vec {
	if s1 == nil || s2 == nil {
		return player.Position
	}
	var posFrom, posTo gmath.Vec
	var foundS1, foundS2 bool
	for i := range s1.Players {
		if s1.Players[i].ID == player.ID {
			converter.StatePosToPlayerPos(&posFrom, s1.Players[i])
			foundS1 = true
			break
		}
	}
	for i := range s2.Players {
		if s2.Players[i].ID == player.ID {
			converter.StatePosToPlayerPos(&posTo, s2.Players[i])
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
