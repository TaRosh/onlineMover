package main

import (
	"image/color"
	"strconv"

	"github.com/TaRosh/online_mover/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Bullet struct {
	game.Bullet
	img *ebiten.Image
	opt *ebiten.DrawImageOptions
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	b.opt.GeoM.Reset()
	b.opt.GeoM.Translate(b.Position.X, b.Position.Y)
	ebitenutil.DebugPrintAt(screen, strconv.Itoa(int(b.ID)), int(b.Position.X+5), int(b.Position.Y+5))
	screen.DrawImage(b.img, b.opt)
}

func NewBullet(id uint32, c color.Color) *Bullet {
	img := ebiten.NewImage(5, 5)
	img.Fill(c)
	b := Bullet{
		Bullet: *game.NewBullet(id),
		img:    img,
		opt:    &ebiten.DrawImageOptions{},
	}
	return &b
}
