package world

import (
	"time"

	"github.com/TaRosh/online_mover/main_client/assets"
	"github.com/TaRosh/online_mover/main_client/player"
	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	shader *ebiten.Shader
	opt    *ebiten.DrawRectShaderOptions
	t      time.Time
}

func (w *World) Draw(screen *ebiten.Image) {
	bounds := screen.Bounds()
	screen.DrawRectShader(bounds.Dx(), bounds.Dy(), w.shader, w.opt)
}

func (w *World) Update(localPlayer *player.Player) {
	w.opt.Uniforms["Time"] = float32(time.Since(w.t).Seconds())
	w.opt.Uniforms["SpeedX"] = float32(10.0 * localPlayer.Velocity.X)
	w.opt.Uniforms["SpeedY"] = float32(10.0*localPlayer.Velocity.Y - 10)
}

func New() *World {
	opt := ebiten.DrawRectShaderOptions{}

	var err error
	t := time.Now()
	opt.Uniforms = map[string]any{
		"Time":   float32(time.Since(t).Seconds()),
		"SpeedX": float32(0.0),
		"SpeedY": float32(-10.0),
	}
	shader, err := assets.GetShader("shaders/starfield.kage")
	if err != nil {
		panic(err)
	}
	w := World{
		shader: shader,
		opt:    &opt,
		t:      t,
	}
	return &w
}
