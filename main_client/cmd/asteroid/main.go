package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	_ "embed"
)

type Game struct{}

var (
	img    *ebiten.Image
	shader *ebiten.Shader
)

//go:embed learn.kage*
var sh []byte

var (
	t   time.Time
	opt ebiten.DrawRectShaderOptions
)

// Draw implements ebiten.Game.
func (g Game) Draw(screen *ebiten.Image) {
	bounds := screen.Bounds()
	screen.DrawRectShader(bounds.Dx(), bounds.Dy(), shader, &opt)
	// o := ebiten.DrawImageOptions{}
	// o.GeoM.Translate(100, 100)
	// screen.DrawImage(img, &o)
}

// Layout implements ebiten.Game.
func (g Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}

// Update implements ebiten.Game.
func (g Game) Update() error {
	opt.Uniforms["Time"] = float32(time.Since(t).Seconds())
	// opt.Uniforms["SpeedX"] = float32(rand.Float32()*2 - 1)
	// opt.Uniforms["SpeedY"] = float32(rand.Float32()*2 - 1)
	return nil
}

func main() {
	var err error
	t = time.Now()
	img = ebiten.NewImage(100, 100)
	opt = ebiten.DrawRectShaderOptions{}
	opt.Uniforms = map[string]any{
		"Time":   t,
		"SpeedX": float32(1.0),
		"SpeedY": float32(0.0),
	}
	shader, err = ebiten.NewShader(sh)
	if err != nil {
		panic(err)
	}
	ebiten.RunGame(Game{})
}
