package game

import (
	"github.com/quasilyte/gmath"
)

type Player struct {
	ID           uint32
	Position     gmath.Vec
	Velocity     gmath.Vec
	Acceleration gmath.Vec
	MaxSpeed     float64
	MaxForce     float64
}

func (p *Player) ApplyForce(force gmath.Vec) {
	p.Acceleration = p.Acceleration.Add(force)
}

func (p *Player) Update() {
	p.Velocity = p.Velocity.Add(p.Acceleration)
	p.Velocity = p.Velocity.ClampLen(p.MaxSpeed)
	p.Position = p.Position.Add(p.Velocity)
	p.Acceleration = p.Acceleration.Mulf(0)
}

func NewPlayer() *Player {
	p := Player{
		MaxSpeed: 5,
		MaxForce: 0.5,
	}
	return &p
}
