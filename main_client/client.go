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
	Width                    int
	Height                   int
	Network                  udp.NetworkClient
	incomingPackets          chan udp.Packet
	lastSnapshotForReconcile *game.Snapshot
	snapshotQueue            []*game.Snapshot
	inputsHistory            []game.Input
	Tick                     uint32
	buf                      []byte
	players                  []*Player
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
	// catch packets
	for {
		select {
		case packet := <-g.incomingPackets:
			// save snapshot in queue
			// save last snapshot for reconcile
			g.processPacket(packet)
		default:
			goto END_NETWORK
		}
	}
END_NETWORK:

	// Reconcile
	if g.lastSnapshotForReconcile != nil {
		g.reapplyPossitionFromSnapshot()
		g.lastSnapshotForReconcile = nil
	}

	// Get inputs
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
	// Apply inputs
	for _, p := range g.players {
		game.ApplyInput(&p.Player, input)
	}
	// Update player with new inputs
	for _, p := range g.players {
		p.Update()
	}
	// Add to input history
	g.inputsHistory = append(g.inputsHistory, input)

	n, err := input.Encode(g.buf)
	if err != nil {
		return nil
	}

	// Send inputs
	err = g.Network.SendInput(g.buf[:n])
	if err != nil {
		panic(err)
	}
	g.Tick += 1

	return nil
}

// Set possiton for player from snapshot
// then reapply inputs from inputs history by tick identifier
// until snapshot last tick input field
func (g *Game) reapplyPossitionFromSnapshot() {
	player := g.players[0]
	player.Position = g.lastSnapshotForReconcile.Players[0].Position
	// player.Velocity = g.lastSnapshotForReconcile.Players[0].Velocity
	inputsAfterSnapshot := g.inputsHistory[:0]
	for _, input := range g.inputsHistory {
		if input.Tick > g.lastSnapshotForReconcile.LastInputTick {
			inputsAfterSnapshot = append(inputsAfterSnapshot, input)
		}
	}
	g.inputsHistory = inputsAfterSnapshot

	for _, input := range g.inputsHistory {
		game.ApplyInput(&player.Player, input)
	}
}

func (g *Game) processPacket(packet udp.Packet) {
	if packet.Type == udp.SnapshotPacket {
		snapshot := game.Snapshot{}
		err := snapshot.Decode(packet.Data)
		if err != nil {
			log.Fatal("precessPacket:packet.Decode", err)
		}
		fmt.Println("Snapshot received:", snapshot)
		// NOT CHANGE SNAPSHOT ITSELF BECOUSE IT'S POINTER!!!
		g.snapshotQueue = append(g.snapshotQueue, &snapshot)
		g.lastSnapshotForReconcile = &snapshot
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
