package main

import (
	"github.com/TaRosh/online_mover/game"
	"github.com/quasilyte/gmath"
)

/*
Convert packet uint position to float vec
for player positios, velocity, etc


*/

func statePosToPlayerPos(pos *gmath.Vec, s game.PlayerState) {
	pos.X = float64(s.X) / 100
	pos.Y = float64(s.Y) / 100
}

// convert state to vector position for player
// position
func stateToPlayer(p *Player, s game.PlayerState) {
	statePosToPlayerPos(&p.Position, s)
}

// convert state to vector position for bullet
// position & velocity
func stateToBullet(b *Bullet, s game.ProjectileState) {
	b.Position.X = float64(s.X) / 100
	b.Position.Y = float64(s.Y) / 100
	b.Velocity.X = float64(s.VX) / 100
	b.Velocity.Y = float64(s.VY) / 100
}
