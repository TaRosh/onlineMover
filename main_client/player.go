package main

import (
	"image/color"
	"math"
	"strconv"

	_ "embed"

	"github.com/TaRosh/online_mover/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed triangle.kage
var triangle []byte

type Player struct {
	game.Player
	img    *ebiten.Image
	shader *ebiten.Shader
	opt    *ebiten.DrawImageOptions
	shOpt  *ebiten.DrawRectShaderOptions
}

func (p *Player) Draw(screen *ebiten.Image) {
	defer p.opt.GeoM.Reset()
	bounds := p.img.Bounds()
	p.opt.GeoM.Translate(-float64(bounds.Dx()/2), -float64(bounds.Dy()/2))
	// angle := p.Velocity.Angle() + math.Pi/2
	p.opt.GeoM.Rotate(float64(p.Rotation) + math.Pi/2)
	p.opt.GeoM.Translate(p.Position.X, p.Position.Y)
	ebitenutil.DebugPrintAt(screen, strconv.Itoa(int(p.ID)), int(p.Position.X+5), int(p.Position.Y+5))
	p.img.DrawRectShader(bounds.Dx(), bounds.Dy(), p.shader, p.shOpt)
	screen.DrawImage(p.img, p.opt)
}

func NewPlayer(id game.PlayerID, c color.Color) (*Player, error) {
	img := ebiten.NewImage(20, 30)
	shader, err := ebiten.NewShader(triangle)
	if err != nil {
		return nil, err
	}
	shaderOpt := &ebiten.DrawRectShaderOptions{}

	p := Player{
		Player: *game.NewPlayer(id),
		img:    img,
		opt:    &ebiten.DrawImageOptions{},
		shader: shader,
		shOpt:  shaderOpt,
	}
	return &p, nil
}
