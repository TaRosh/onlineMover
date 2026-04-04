package main

import (
	"fmt"
	"os"
	"time"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/network"
	"github.com/joho/godotenv"
	"github.com/quasilyte/gmath"
)

const tickTime = time.Second / 20

type World struct {
	players      map[game.PlayerID]*game.Player
	bulletNextID uint32
	bullets      map[uint32]*game.Bullet
	inputsQueue  chan game.Input
	// for now we have one event: receive new player id
	incomingEvents      chan game.Event
	lastPlayerSnapshots map[game.PlayerID]game.PlayerState
	lastBulletSnapshots map[uint32]game.ProjectileState
	currentSnapshot     *game.Snapshot
	tick                uint32
	network             network.NetworkServer
	buf                 []byte
}

func (w *World) Init() {
	w.players = make(map[game.PlayerID]*game.Player)
	w.bullets = make(map[uint32]*game.Bullet)
	w.lastPlayerSnapshots = make(map[game.PlayerID]game.PlayerState)
	w.lastBulletSnapshots = map[uint32]game.ProjectileState{}
	w.inputsQueue = make(chan game.Input, 1024)
	w.incomingEvents = make(chan game.Event, 1024)
	w.currentSnapshot = new(game.Snapshot)
	w.tick = 0
	network, err := network.NewServer(os.Getenv("S_HOST"), os.Getenv("S_PORT"), []byte(os.Getenv("ENC_KEY")))
	if err != nil {
		// can't create udp server
		panic(err)
	}
	w.network = network
	w.buf = make([]byte, 2048)
}

func (w *World) processEvent(event game.Event) {
	switch event.Type {
	case game.EventConnection:
		player := game.NewPlayer(event.ID)
		w.players[player.ID] = player
		// TODO: notify players about new player
		w.sendFullSnapshot()
	case game.EventNoAnswerFromClient:
		w.network.DeletePlayer(event.ID)
		delete(w.players, event.ID)
	}
}

// world.spawnBullet(p)
func (w *World) spawnBullet(p *game.Player) {
	b := game.NewBullet(w.bulletNextID)
	b.Position = p.Position
	direction := p.Rotation
	b.Velocity = gmath.Vec{X: direction.Cos(), Y: direction.Sin()}.Mulf(5)
	// b.Velocity.X = 10

	w.bullets[w.bulletNextID] = b
	w.bulletNextID += 1
}

func main() {
	ticker := time.Tick(tickTime)
	world := World{}
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	world.Init()
	// world.players = append(world.players, game.NewPlayer())
	// now we don't know amount of players
	// world.snapshot.Players = make([]game.PlayerState, len(world.players))
	go world.network.Receive(world.inputsQueue, world.incomingEvents)
	fmt.Println("Server ready")
	for range ticker {
		// connect new players first
		for {
			select {
			case event := <-world.incomingEvents:
				// save inputs to game
				world.processEvent(event)
			default:
				goto END_EVENTS
			}
		}
	END_EVENTS:
		// update world on network inputs result
		for {
			select {
			case input := <-world.inputsQueue:
				p := world.players[input.ID]
				game.ApplyInput(p, input)
				if input.Buttons&game.InputShoot != 0 {
					if world.tick-p.LastFireTick > 10 {
						world.spawnBullet(p)
						p.LastFireTick = world.tick
					}
				}
				// TODO: think about this can be unsync in ticks
				// from clients?
				world.currentSnapshot.LastInputTick = input.Tick
			default:
				goto END_APPLY_INPUT
			}
		}
	END_APPLY_INPUT:
		world.update()
		// send snapshot on network layer
		// time.Sleep(time.Second)
		world.sendDeltaSnapshot()
		world.tick += 1
		world.network.CheckTimeouts(world.incomingEvents)
	}
}

// TODO:
// add full snapshot
func (w *World) sendFullSnapshot() {
	if len(w.players) < 1 {
		return
	}

	w.currentSnapshot.Tick = w.tick
	w.currentSnapshot.Players = nil
	w.currentSnapshot.Projectiles = nil

	playerState := game.PlayerState{}
	for _, p := range w.players {
		playerState.ID = p.ID
		playerState.X = int32(p.Position.X * 100)
		playerState.Y = int32(p.Position.Y * 100)
		playerState.VX = int16(p.Velocity.X * 100)
		playerState.VY = int16(p.Velocity.Y * 100)
		w.currentSnapshot.Players = append(w.currentSnapshot.Players, playerState)
		w.lastPlayerSnapshots[p.ID] = playerState
	}
	projectileState := game.ProjectileState{}
	for _, b := range w.bullets {
		projectileState.ID = b.ID
		projectileState.X = int32(b.Position.X * 100)
		projectileState.Y = int32(b.Position.Y * 100)
		projectileState.VX = int16(b.Velocity.X * 100)
		projectileState.VY = int16(b.Velocity.Y * 100)
		w.currentSnapshot.Projectiles = append(w.currentSnapshot.Projectiles, projectileState)
		w.lastBulletSnapshots[b.ID] = projectileState
	}
	n, err := w.currentSnapshot.Encode(w.buf)
	if err != nil {
		panic(err)
	}
	fmt.Println("snapshot bytes:", len(w.buf[:n]))
	fmt.Printf("SENT snapshot: %+v\n", w.currentSnapshot)
	w.broadcastSnapshot(w.buf[:n])

	// last snapshot needed to compare entity changes
}

// and delta snapshot

func (w *World) sendDeltaSnapshot() {
	if len(w.players) < 1 {
		if w.tick > 1000 {
			panic("panic")
		}
		return
	}

	w.currentSnapshot.Tick = w.tick
	w.currentSnapshot.Players = nil
	w.currentSnapshot.Projectiles = nil

	playerState := game.PlayerState{}
	for _, p := range w.players {
		playerState.ID = p.ID
		playerState.X = int32(p.Position.X * 100)
		playerState.Y = int32(p.Position.Y * 100)
		playerState.VX = int16(p.Velocity.X * 100)
		playerState.VY = int16(p.Velocity.Y * 100)
		// TODO: fill it
		lastSnapshot, exist := w.lastPlayerSnapshots[p.ID]
		// if !exist {
		// 	w.currentSnapshot.Players = append(w.currentSnapshot.Players, playerState)
		// } else if lastSnapshot != playerState {
		// 	w.currentSnapshot.Players = append(w.currentSnapshot.Players, playerState)
		// }
		if !exist || lastSnapshot != playerState {
			w.currentSnapshot.Players = append(w.currentSnapshot.Players, playerState)
		}
		w.lastPlayerSnapshots[p.ID] = playerState
	}
	projectileState := game.ProjectileState{}
	for _, b := range w.bullets {
		projectileState.ID = b.ID
		projectileState.X = int32(b.Position.X * 100)
		projectileState.Y = int32(b.Position.Y * 100)
		projectileState.VX = int16(b.Velocity.X * 100)
		projectileState.VY = int16(b.Velocity.Y * 100)
		lastSnapshot, exist := w.lastBulletSnapshots[b.ID]
		if !exist || lastSnapshot != projectileState {
			w.currentSnapshot.Projectiles = append(w.currentSnapshot.Projectiles, projectileState)
		}
		w.lastBulletSnapshots[b.ID] = projectileState
	}
	n, err := w.currentSnapshot.Encode(w.buf)
	if err != nil {
		panic(err)
	}
	fmt.Println("snapshot bytes:", len(w.buf[:n]))
	fmt.Printf("SENT snapshot: %+v\n", w.currentSnapshot)
	w.broadcastSnapshot(w.buf[:n])

	// last snapshot needed to compare entity changes
}

func (w *World) broadcastSnapshot(snapshot []byte) {
	for id := range w.players {
		err := w.network.SendSnapshot(id, snapshot)
		if err != nil {
			// TODO: think about send snapshot error
			panic(err)
		}
	}
}

func (w *World) update() {
	for _, p := range w.players {
		p.Update()
	}
	for _, b := range w.bullets {
		b.Update()
	}
}
