package main

import (
	"fmt"
	"os"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/network"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joho/godotenv"
)

func (g *Game) Init() {
	g.Width, g.Height = ebiten.WindowSize()
	g.players = make(map[game.PlayerID]*Player)
	g.projectiles = make(map[uint32]*Bullet)
	network, err := network.NewClient(os.Getenv("CL_HOST"), os.Getenv("CL_PORT"), []byte(os.Getenv("ENC_KEY")))
	if err != nil {
		// TODO: think about reconnect to server
		panic(err)
	}
	g.Network = network
	g.snapshotQueue = make(chan game.Snapshot, 1024)
	g.maxSnapshot = 64
	g.events = make(chan game.Event, 1024)
	g.Tick = 0
	// ~ 100ms
	// render in past about 2 ticks
	g.interpolationDelay = 2
	g.buf = make([]byte, 1024)
}

func main() {
	var err error
	g := Game{}
	err = godotenv.Load()
	if err != nil {
		panic(err)
	}
	g.Init()
	// g.snapshotQueue = make(chan *game.Snapshot, 1024)
	// send request about connect player to server
	go g.Network.Receive(g.snapshotQueue, g.events)
	g.connectPlayerToServer()

	// change tick on client side?
	ebiten.SetTPS(20)
	err = ebiten.RunGame(&g)
	if err != nil {
		fmt.Println("RunGame:", err)
	}
}
