package entities

import (
	"github.com/TaRosh/online_mover/game/layers"
	quadtree "github.com/TaRosh/online_mover/quad_tree"
	"github.com/quasilyte/gmath"
)

type Bullet struct {
	ID uint32

	Width  float64
	Height float64

	Deleted bool

	Position     gmath.Vec
	Velocity     gmath.Vec
	Acceleration gmath.Vec

	Rotation gmath.Rad

	MaxSpeed float64
	MaxForce float64
}

func (b *Bullet) Update() {
	b.Velocity = b.Velocity.Add(b.Acceleration)
	b.Velocity = b.Velocity.ClampLen(b.MaxSpeed)
	if b.Velocity.Len() > 0.001 {
		b.Rotation = b.Velocity.Angle()
	}
	b.Position = b.Position.Add(b.Velocity)
	b.Acceleration = b.Acceleration.Mulf(0)
}

func (b *Bullet) Id() uint32 {
	return b.ID
}

func (b *Bullet) Bounds() quadtree.Rect[float64] {
	rect := quadtree.NewRect(b.Position.X-b.Width/2, b.Position.Y-b.Height/2,
		b.Position.X+b.Width/2, b.Position.Y+b.Height/2)
	return rect
}

func (b *Bullet) Layer() quadtree.Layer {
	return layers.LayerBullet
}

func (b *Bullet) Mask() quadtree.Layer {
	return layers.LayerAsteroid
}

func NewBullet(id uint32, w, h float64) *Bullet {
	b := Bullet{
		ID:       id,
		Width:    w,
		Height:   h,
		MaxSpeed: 5,
		MaxForce: 0.5,
	}
	return &b
}
