package game

import (
	"github.com/quasilyte/gmath"
)

type PlayerID uint32

type Player struct {
	ID           PlayerID
	Position     gmath.Vec
	Velocity     gmath.Vec
	Acceleration gmath.Vec
	LastFireTick uint32
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

func NewPlayer(id PlayerID) *Player {
	p := Player{
		ID:       id,
		MaxSpeed: 5,
		MaxForce: 0.5,
	}
	return &p
}
