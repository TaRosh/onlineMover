package game

import (
	engine "github.com/TaRosh/online_mover/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func stateInput(g *Game) stateFn {
	g.buttons = 0
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.buttons |= engine.InputUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.buttons |= engine.InputDown
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.buttons |= engine.InputLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.buttons |= engine.InputRight
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.buttons |= engine.InputShoot
	}
	g.input.ID = g.localPlayer.ID
	g.input.Tick = g.Tick
	g.input.Buttons = g.buttons
	return nil
}
