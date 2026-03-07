package main

import (
	"fmt"
	"image/color"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/udp"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Width   int
	Height  int
	Network udp.NetworkClient
	Tick    uint32
	buf     []byte
	players []*game.Player
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.players[0].Draw(screen)
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Width, g.Height
}

var (
	buttons uint8
	input   game.Input
)

// send input by tick
func (g *Game) Update() error {
	buttons = 0
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		buttons |= game.InputUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		buttons |= game.InputDown
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		buttons |= game.InputLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		buttons |= game.InputRight
	}
	input.Tick = g.Tick
	input.Buttons = buttons
	game.ApplyInput(g.players[0], input)
	// n, err := input.Encode(g.buf)
	// if err != nil {
	// 	return nil
	// }

	// err = g.Network.SendInput(g.buf[:n])
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("sending input: %+v\n", input)
	g.Tick += 1

	return nil
}

func main() {
	var err error
	g := Game{}
	g.Width, g.Height = ebiten.WindowSize()
	g.players = append(g.players, game.NewPlayer(color.White))
	g.Network, err = udp.NewClient("localhost", "9000")
	g.buf = make([]byte, 1024)
	if err != nil {
		panic(err)
	}
	// change tick on client side?
	ebiten.SetTPS(20)
	err = ebiten.RunGame(&g)
	if err != nil {
		fmt.Println("RunGame:", err)
	}
}
