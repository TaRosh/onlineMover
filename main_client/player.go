package main

import (
	"image/color"
	"strconv"

	"github.com/TaRosh/online_mover/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Player struct {
	game.Player
	img *ebiten.Image
	opt *ebiten.DrawImageOptions
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.opt.GeoM.Reset()
	p.opt.GeoM.Translate(p.Position.X, p.Position.Y)
	ebitenutil.DebugPrintAt(screen, strconv.Itoa(int(p.ID)), int(p.Position.X+5), int(p.Position.Y+5))
	screen.DrawImage(p.img, p.opt)
}

func NewPlayer(id game.PlayerID, c color.Color) *Player {
	img := ebiten.NewImage(10, 10)
	img.Fill(c)
	p := Player{
		Player: *game.NewPlayer(id),
		img:    img,
		opt:    &ebiten.DrawImageOptions{},
	}
	return &p
}
