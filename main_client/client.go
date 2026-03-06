package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/udp"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/gmath"
)

type Game struct {
	Width   int
	Height  int
	Network udp.NetworkClient
	Tick    uint32
	buf     []byte
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Width, g.Height
}

// send input by tick
func (g *Game) Update() error {
	x := rand.Float64()
	y := rand.Float64()
	vel := gmath.Vec{X: x, Y: y}
	input := game.Input{Tick: g.Tick, Vel: vel}
	n, err := input.Encode(g.buf)
	if err != nil {
		return nil
	}

	err = g.Network.SendInput(g.buf[:n])
	if err != nil {
		panic(err)
	}
	fmt.Printf("sending input: %+v\n", input)
	g.Tick += 1

	return nil
}

func main() {
	var err error
	g := Game{}
	g.Width, g.Height = ebiten.WindowSize()
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
