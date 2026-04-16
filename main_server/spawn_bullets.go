package main

import (
	"github.com/TaRosh/online_mover/game/entities"
	"github.com/quasilyte/gmath"
)

func (w *World) spawnBullet(p *entities.Player) {
	b := entities.NewBullet(w.bulletNextID, 32, 32)
	b.Position = p.Position
	direction := p.Rotation
	b.Velocity = gmath.Vec{X: direction.Cos(), Y: direction.Sin()}.Mulf(5)
	// b.Velocity.X = 10

	w.bullets[w.bulletNextID] = b
	w.bulletNextID += 1
}
