package asteroid

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func setShaderOpt(r float32) (*ebiten.DrawRectShaderOptions, error) {
	shaderOpt := ebiten.DrawRectShaderOptions{}
	// var err error
	shaderOpt.Uniforms = map[string]any{
		"Time":   time.Now(),
		"Radius": r,
		"Color":  []float32{1.0, 1.0, 0},
	}
	return &shaderOpt, nil
}
