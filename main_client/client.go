package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/udp"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Width           int
	Height          int
	Network         udp.NetworkClient
	incomingPackets chan udp.Packet
	snapshotQueue   []*game.Snapshot
	Tick            uint32
	buf             []byte
	players         []*Player
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.players {
		p.Draw(screen)
	}
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
	for _, p := range g.players {
		game.ApplyInput(&p.Player, input)
	}
	for _, p := range g.players {
		p.Update()
	}

	n, err := input.Encode(g.buf)
	if err != nil {
		return nil
	}

	err = g.Network.SendInput(g.buf[:n])
	if err != nil {
		panic(err)
	}
	// catch packets
	for {
		select {
		case packet := <-g.incomingPackets:
			// save snapshot in queue
			g.processPacket(packet)
		default:
			goto END_NETWORK
		}
	}
END_NETWORK:
	// Apply snapshot
	// fmt.Printf("sending input: %+v\n", input)
	g.Tick += 1

	return nil
}

func (g *Game) processPacket(packet udp.Packet) {
	fmt.Printf("Packet receive: %+v\n", packet)
	if packet.Type == udp.SnapshotPacket {
		snapshot := game.Snapshot{}
		err := snapshot.Decode(packet.Data)
		if err != nil {
			log.Fatal("precessPacket:packet.Decode", err)
		}
		g.snapshotQueue = append(g.snapshotQueue, &snapshot)
	}
}

func main() {
	var err error
	g := Game{}
	g.Width, g.Height = ebiten.WindowSize()
	g.players = append(g.players, NewPlayer(color.White))
	g.Network, err = udp.NewClient("localhost", "9000")
	g.incomingPackets = make(chan udp.Packet, 1024)
	// g.snapshotQueue = make(chan *game.Snapshot, 1024)
	go g.Network.Receive(g.incomingPackets)
	// now we need receive snapshot
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
