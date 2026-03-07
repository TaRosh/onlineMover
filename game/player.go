package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/gmath"
)

type Player struct {
	ID           uint32
	Position     gmath.Vec
	Velocity     gmath.Vec
	Acceleration gmath.Vec
	MaxSpeed     float64
	MaxForce     float64
	img          *ebiten.Image
	opt          *ebiten.DrawImageOptions
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

func (p *Player) Draw(screen *ebiten.Image) {
	p.opt.GeoM.Reset()
	p.opt.GeoM.Translate(p.Position.X, p.Position.Y)
	screen.DrawImage(p.img, p.opt)
}

func NewPlayer(c color.Color) *Player {
	img := ebiten.NewImage(10, 10)
	img.Fill(c)
	p := Player{
		img:      img,
		opt:      &ebiten.DrawImageOptions{},
		MaxSpeed: 5,
		MaxForce: 0.5,
	}
	return &p
}
