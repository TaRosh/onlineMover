package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var t time.Time = time.Now()

func (g *Game) Draw(screen *ebiten.Image) {
	// img := ebiten.NewImage(int(g.camera.Width), int(g.camera.Height))
	// img.Fill(color.RGBA{0xff, 0xff, 0xff, 0x50})
	// screen.DrawImage(img, nil)
	g.world.Draw(screen)
	for _, p := range g.players {
		p.Draw(screen, g.camera)
	}
	for _, b := range g.projectiles {
		b.Draw(screen, g.camera)
	}
	for _, a := range g.asteroids {
		a.Draw(screen, g.camera)
	}
}
