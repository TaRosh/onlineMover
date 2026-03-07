package main

import (
	"image/color"

	"github.com/TaRosh/online_mover/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	game.Player
	img *ebiten.Image
	opt *ebiten.DrawImageOptions
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
		Player: *game.NewPlayer(),
		img:    img,
		opt:    &ebiten.DrawImageOptions{},
	}
	return &p
}
