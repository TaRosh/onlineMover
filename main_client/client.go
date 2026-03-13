package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/network"
	"github.com/TaRosh/online_mover/udp"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Width                    int
	Height                   int
	Network                  network.NetworkClient
	lastSnapshotForReconcile *game.Snapshot
	snapshotQueue            []*game.Snapshot
	inputsHistory            []game.Input

	Tick uint32
	buf  []byte
	// separate becouse we want recocilate
	localPlayer *Player
	players     map[uint32]*Player
	// debugPlayer *Player
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.players {
		p.Draw(screen)
	}
	// g.debugPlayer.Draw(screen)
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

		for _, player := range g.lastSnapshotForReconcile.Players {
			_, exist := g.players[player.ID]
			if !exist {
				g.players[player.ID] = NewPlayer(player.ID, color.RGBA{0xff, 0, 0, 0xff})
			}

			if player.ID == g.localPlayer.ID {
				g.reapplyPossitionFromSnapshot(player)
			} else {
				g.players[player.ID].Position = player.Position
			}
		}

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
	input.ID = g.localPlayer.ID
	input.Tick = g.Tick
	input.Buttons = buttons
	// Apply inputs
	// for _, p := range g.players {
	game.ApplyInput(&g.localPlayer.Player, input)
	// }
	// Update player with new inputs
	g.localPlayer.Update()
	// for _, p := range g.players {
	// 	p.Update()
	// }
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
func (g *Game) reapplyPossitionFromSnapshot(player game.PlayerState) {
	g.localPlayer.Position = player.Position
	// player.Velocity = g.lastSnapshotForReconcile.Players[0].Velocity
	inputsAfterSnapshot := g.inputsHistory[:0]
	for _, input := range g.inputsHistory {
		if input.Tick > g.lastSnapshotForReconcile.LastInputTick {
			inputsAfterSnapshot = append(inputsAfterSnapshot, input)
		}
	}
	g.inputsHistory = inputsAfterSnapshot

	for _, input := range g.inputsHistory {
		game.ApplyInput(&g.localPlayer.Player, input)
	}
	// g.debugPlayer.Position = *&player.Position
	// fmt.Println("DEBUG PLAYER POS", g.debugPlayer.Position)
}

func (g *Game) processPacket(packet udp.Packet) {
	// first we need get our connectian accept from server
	// server give our player ID ! ! ! ( MATER )
	if g.localPlayer == nil && packet.Type != udp.AcceptPacket {
		// we don't have local player and still don't get
		// answer for our playerID request
		// maybe resend after some time?
		return
	}

	if packet.Type == udp.SnapshotPacket {
		snapshot := game.Snapshot{}
		err := snapshot.Decode(packet.Data)
		if err != nil {
			log.Fatal("precessPacket:packet.Decode", err)
		}
		// fmt.Println("Snapshot received:", snapshot)
		// NOT CHANGE SNAPSHOT ITSELF BECOUSE IT'S POINTER!!!
		g.snapshotQueue = append(g.snapshotQueue, &snapshot)
		g.lastSnapshotForReconcile = &snapshot
	}
}

func (g *Game) connectPlayerToServer() {
	tries := 4
	// time after we resend connection request
	resendAfter := time.Tick(time.Millisecond * 50)
	// for range tries {
TRY_GET_PLAYER_ID_AGAIN:
	if tries == 0 {
		panic("could'n get player id from server. I'm done! Shut Done")
	}
	err := g.Network.SendPlayerConnectionRequest()
	// another try when network problem
	if err != nil {
		tries--
		time.Sleep(time.Millisecond * 50)
		goto TRY_GET_PLAYER_ID_AGAIN
	}
	// another try when not answer from server
	for {
		select {
		case p := <-g.incomingPackets:
			if p.Type != udp.AcceptPacket {
				break
			}
			id := game.PlayerIDPacket{}
			err := id.Decode(p.Data)
			if err != nil {
				// can't decode my id
				panic(err)
			}
			fmt.Println("MY ID", id.ID)
			player := NewPlayer(id.ID, color.White)
			g.localPlayer = player
			g.players[id.ID] = player
			return
		case <-resendAfter:
			tries--
			goto TRY_GET_PLAYER_ID_AGAIN
		}
	}
}

func main() {
	var err error
	g := Game{}
	g.Width, g.Height = ebiten.WindowSize()
	g.players = make(map[uint32]*Player)
	// g.debugPlayer = NewPlayer(color.RGBA{0xff, 0, 0, 0xff})
	g.Network, err = udp.NewClient("localhost", "9000")
	g.incomingPackets = make(chan udp.Packet, 1024)
	// g.snapshotQueue = make(chan *game.Snapshot, 1024)
	// send request about connect player to server
	go g.Network.Receive(g.incomingPackets)
	g.connectPlayerToServer()

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
