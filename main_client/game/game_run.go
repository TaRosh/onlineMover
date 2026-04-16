package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Run() {
	go g.Network.Receive(g.snapshotQueue, g.events)
	// TODO: uncomment line bellow
	g.connectPlayerToServer()

	// change tick on client side?
	ebiten.SetTPS(20)
	err := ebiten.RunGame(g)
	if err != nil {
		fmt.Println("RunGame:", err)
	}
}
