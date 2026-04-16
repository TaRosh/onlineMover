package entities

import (
	"github.com/TaRosh/online_mover/game/layers"
	quadtree "github.com/TaRosh/online_mover/quad_tree"
	"github.com/quasilyte/gmath"
)

type Asteroid struct {
	ID uint32

	Width  float64
	Height float64

	Deleted bool

	Position     gmath.Vec
	Velocity     gmath.Vec
	Acceleration gmath.Vec

	Rotation  gmath.Rad
	RotationV gmath.Rad

	MaxSpeed float64
	MaxForce float64
}

func (a *Asteroid) ApplyForce(force gmath.Vec) {
	a.Acceleration = a.Acceleration.Add(force)
}

func (a *Asteroid) Update() {
	a.Velocity = a.Velocity.Add(a.Acceleration)
	a.Velocity = a.Velocity.ClampLen(a.MaxSpeed)
	// a.RotationV += gmath.Rad(rand.Float32()*2 - 1)
	// a.RotationV = gmath.Clamp(a.RotationV, -0.1, 0.1)
	// a.Rotation = a.Rotation + a.RotationV
	a.Position = a.Position.Add(a.Velocity)
	a.Acceleration = a.Acceleration.Mulf(0)
}

func (a *Asteroid) Id() uint32 {
	return a.ID
}

func (a *Asteroid) Bounds() quadtree.Rect[float64] {
	rect := quadtree.NewRect(a.Position.X-a.Width/2, a.Position.Y-a.Height/2,
		a.Position.X+a.Width/2, a.Position.Y+a.Height/2)
	return rect
}

func (a *Asteroid) Layer() quadtree.Layer {
	return layers.LayerAsteroid
}

func (a *Asteroid) Mask() quadtree.Layer {
	return layers.LayerPlayer | layers.LayerBullet
}

func NewAsteroid(id uint32, w, h float64) *Asteroid {
	a := Asteroid{
		ID:       id,
		Width:    w,
		Height:   h,
		MaxSpeed: 5,
		MaxForce: 5,
	}
	return &a
}
