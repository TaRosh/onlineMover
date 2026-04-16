package player

import (
	"github.com/TaRosh/online_mover/main_client/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

func setShaderOpt() (*ebiten.DrawRectShaderOptions, error) {
	shaderOpt := ebiten.DrawRectShaderOptions{}
	var err error
	shaderOpt.Images[0], err = assets.GetImage("img/ship/ship_full.png")
	if err != nil {
		return nil, err
	}
	shaderOpt.Images[1], err = assets.GetImage("img/ship/ship_light_damage.png")
	if err != nil {
		return nil, err
	}
	shaderOpt.Images[2], err = assets.GetImage("img/ship/ship_damaged.png")
	if err != nil {
		return nil, err
	}
	shaderOpt.Images[3], err = assets.GetImage("img/ship/ship_very_damaged.png")
	if err != nil {
		return nil, err
	}
	shaderOpt.Uniforms = map[string]any{
		"HP": 1.0,
	}
	return &shaderOpt, nil
}
