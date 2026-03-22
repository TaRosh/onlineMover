package game

import "github.com/quasilyte/gmath"

type Bullet struct {
	ID           uint32
	Position     gmath.Vec
	Velocity     gmath.Vec
	Acceleration gmath.Vec
	MaxSpeed     float64
	MaxForce     float64
}

func (b *Bullet) Update() {
	b.Velocity = b.Velocity.Add(b.Acceleration)
	b.Velocity = b.Velocity.ClampLen(b.MaxSpeed)
	b.Position = b.Position.Add(b.Velocity)
	b.Acceleration = b.Acceleration.Mulf(0)
}

func NewBullet(id uint32) *Bullet {
	b := Bullet{
		ID:       id,
		MaxSpeed: 5,
		MaxForce: 0.5,
	}
	return &b
}
