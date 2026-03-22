package main

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/network"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joho/godotenv"
	"github.com/quasilyte/gmath"
)

type Game struct {
	Width   int
	Height  int
	Network network.NetworkClient

	lastServerTick           uint32
	lastSnapshotForReconcile *game.Snapshot
	interpolationDelay       uint32

	// snapshotQueue []*game.Snapshot
	snapshotQueue chan game.Snapshot
	// length of snapshot history
	maxSnapshot    int
	snapshotBuffer []*game.Snapshot
	events         chan game.Event
	inputsHistory  []game.Input

	Tick uint32
	buf  []byte
	// separate becouse we want recocilate
	localPlayer *Player
	players     map[game.PlayerID]*Player
	// debugPlayer *Player
	projectiles map[uint32]*Bullet
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.players {
		p.Draw(screen)
	}
	for _, b := range g.projectiles {
		b.Draw(screen)
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

func (g *Game) playerInterpolation(player *Player) gmath.Vec {
	renderTick := g.lastServerTick - g.interpolationDelay
	var s1, s2 *game.Snapshot
	for i := 0; i < len(g.snapshotBuffer)-1; i++ {
		if g.snapshotBuffer[i].Tick <= renderTick &&
			g.snapshotBuffer[i+1].Tick >= renderTick {
			s1 = g.snapshotBuffer[i]
			s2 = g.snapshotBuffer[i+1]
			break
		}
	}
	if s1 == nil || s2 == nil {
		return player.Position
	}
	alpha := float64((renderTick - s1.Tick)) / float64((s2.Tick - s1.Tick))
	var posFrom, posTo gmath.Vec
	var foundS1, foundS2 bool
	for i := range s1.Players {
		if s1.Players[i].ID == player.ID {
			statePosToPlayerPos(&posFrom, s1.Players[i])
			foundS1 = true
			break
		}
	}
	for i := range s2.Players {
		if s2.Players[i].ID == player.ID {
			statePosToPlayerPos(&posTo, s2.Players[i])
			foundS2 = true
			break
		}
	}
	if !foundS1 || !foundS2 {
		return player.Position
	}
	newPos := gmath.Vec{
		X: gmath.Lerp(posFrom.X, posTo.X, alpha),
		Y: gmath.Lerp(posFrom.Y, posTo.Y, alpha),
	}
	return newPos
}

func (g *Game) updateProjectiles() {
	for _, b := range g.projectiles {
		b.Update()
	}
}

func statePosToPlayerPos(pos *gmath.Vec, s game.PlayerState) {
	pos.X = float64(s.X) / 100
	pos.Y = float64(s.Y) / 100
}

// convert state to vector position for player
func stateToPlayer(p *Player, s game.PlayerState) {
	statePosToPlayerPos(&p.Position, s)
}

// convert state to vector position for bullet
func stateToBullet(b *Bullet, s game.ProjectileState) {
	b.Position.X = float64(s.X) / 100
	b.Position.Y = float64(s.Y) / 100
	b.Velocity.X = float64(s.VX) / 100
	b.Velocity.Y = float64(s.VY) / 100
}

func (g *Game) processSnapshot(snapshot *game.Snapshot) {
	for _, player := range snapshot.Players {
		_, exist := g.players[player.ID]
		if !exist {
			p := NewPlayer(player.ID, color.RGBA{0xff, 0, 0, 0xff})
			stateToPlayer(p, player)
			g.players[player.ID] = p
			fmt.Println("HERE CREATE PLAYER ")
		}
	}
}

// send input by tick
func (g *Game) Update() error {
	// catch packets
	for {
		select {
		case snapshot := <-g.snapshotQueue:
			// add only newer snapshots
			fmt.Printf("%+v\n", snapshot)
			if snapshot.Tick > g.lastServerTick {
				if len(g.snapshotBuffer) > g.maxSnapshot {
					g.snapshotBuffer = g.snapshotBuffer[1:]
				}
				g.snapshotBuffer = append(g.snapshotBuffer, &snapshot)
				g.lastServerTick = snapshot.Tick
			}
			g.lastSnapshotForReconcile = &snapshot
			g.processSnapshot(&snapshot)
		default:
			goto END_NETWORK
		}
	}
END_NETWORK:

	// Reconcile
	if g.lastSnapshotForReconcile != nil {

		for _, player := range g.lastSnapshotForReconcile.Players {
			if player.ID == g.localPlayer.ID {
				g.reapplyPossitionFromSnapshot(player)
			} else {
				// TODO: add players deletion
				g.players[player.ID].Position = g.playerInterpolation(g.players[player.ID])
			}
		}
		// get projectiles from snapshot
		for _, bullet := range g.lastSnapshotForReconcile.Projectiles {

			b, exist := g.projectiles[bullet.ID]
			if !exist {
				b = NewBullet(bullet.ID, color.RGBA{0, 0, 0xff, 0xff})
				g.projectiles[b.ID] = b
			}
			stateToBullet(b, bullet)
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
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		buttons |= game.InputShoot
	}
	// TODO: add rotation
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
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
	for _, b := range g.projectiles {
		b.Update()
	}
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
	newPos := gmath.Vec{}
	statePosToPlayerPos(&newPos, player)
	g.localPlayer.Position = newPos
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

func (g *Game) connectPlayerToServer() {
	tries := 4
	// time after we resend connection request
	resendAfter := time.Tick(time.Second * 2)
	// for range tries {
TRY_GET_PLAYER_ID_AGAIN:
	if tries == 0 {
		panic("could'n get player id from server. I'm done! Shut Done")
	}
	err := g.Network.SendPlayerConnectionRequest()
	// another try when network problem
	if err != nil {
		tries--
		time.Sleep(time.Millisecond * 300)
		goto TRY_GET_PLAYER_ID_AGAIN
	}
	// another try when no answer from server
	for {
		select {
		case event := <-g.events:
			if event.Type == game.EventConnection {
				id := event.ID
				fmt.Println("MY ID", id)
				player := NewPlayer(id, color.White)
				g.localPlayer = player
				g.players[id] = player
			}
			return
		case <-resendAfter:
			tries--
			goto TRY_GET_PLAYER_ID_AGAIN
		}
	}
}

func (g *Game) Init() {
	g.Width, g.Height = ebiten.WindowSize()
	g.players = make(map[game.PlayerID]*Player)
	g.projectiles = make(map[uint32]*Bullet)
	network, err := network.NewClient(os.Getenv("CL_HOST"), os.Getenv("CL_PORT"))
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
