package main

import "github.com/hajimehoshi/ebiten/v2"

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.players {
		p.Draw(screen)
	}
	for _, b := range g.projectiles {
		b.Draw(screen)
	}
	// g.debugPlayer.Draw(screen)
}
