package entities

import (
	"github.com/TaRosh/online_mover/game/layers"
	quadtree "github.com/TaRosh/online_mover/quad_tree"
	"github.com/quasilyte/gmath"
)

type PlayerID uint32

type Player struct {
	ID           PlayerID
	Position     gmath.Vec
	Velocity     gmath.Vec
	Acceleration gmath.Vec
	Rotation     gmath.Rad
	Health       int

	Width  float64
	Height float64

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
	if p.Velocity.Len() > 0.001 {
		p.Rotation = p.Velocity.Angle()
	}
	p.Position = p.Position.Add(p.Velocity)
	p.Acceleration = p.Acceleration.Mulf(0)
}

func (p *Player) Id() uint32 {
	return uint32(p.ID)
}

func (p *Player) Layer() quadtree.Layer {
	return layers.LayerPlayer
}

func (p *Player) Mask() quadtree.Layer {
	return layers.LayerAsteroid
}

func (p *Player) Bounds() quadtree.Rect[float64] {
	rect := quadtree.NewRect(p.Position.X-p.Width/2, p.Position.Y-p.Height/2,
		p.Position.X+p.Width/2, p.Position.Y+p.Height/2)
	return rect
}

func NewPlayer(id PlayerID, w, h float64) *Player {
	p := Player{
		ID:       id,
		Width:    w,
		Height:   h,
		MaxSpeed: 20,
		MaxForce: 5,
		Health:   100,
	}
	return &p
}
