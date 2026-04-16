package converter

import (
	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/main_client/asteroid"
	"github.com/TaRosh/online_mover/main_client/bullet"
	"github.com/TaRosh/online_mover/main_client/player"
	"github.com/quasilyte/gmath"
)

/*
Convert packet uint position to float vec
for player positios, velocity, etc


*/

func StatePosToPlayerPos(pos *gmath.Vec, s game.PlayerState) {
	pos.X = float64(s.X) / 100
	pos.Y = float64(s.Y) / 100
}

func StateAngleToPlayerRotation(rotation *gmath.Rad, s game.PlayerState) {
	angle := gmath.Rad(game.Uint16ToRad(s.Angle))
	*rotation = angle
}

// convert state to vector position for player
// position
func StateToPlayer(p *player.Player, s game.PlayerState) {
	StatePosToPlayerPos(&p.Position, s)
	StateAngleToPlayerRotation(&p.Rotation, s)
}

// convert state to vector position for bullet
// position & velocity
func StateToBullet(b *bullet.Bullet, s game.ProjectileState) {
	b.Position.X = float64(s.X) / 100
	b.Position.Y = float64(s.Y) / 100
	b.Velocity.X = float64(s.VX) / 100
	b.Velocity.Y = float64(s.VY) / 100
}

func StateToAsteroid(b *asteroid.Asteroid, s game.ProjectileState) {
	b.Position.X = float64(s.X) / 100
	b.Position.Y = float64(s.Y) / 100
	b.Velocity.X = float64(s.VX) / 100
	b.Velocity.Y = float64(s.VY) / 100
}
