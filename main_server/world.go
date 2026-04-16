package main

import (
	"os"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/game/entities"
	"github.com/TaRosh/online_mover/network"
	quadtree "github.com/TaRosh/online_mover/quad_tree"
	"github.com/quasilyte/gmath"
)

type World struct {
	Width  int
	Height int

	// position of this point
	playerSpawnPoint gmath.Vec
	// width height
	asteroidSpawnArea gmath.Vec

	qtree *quadtree.QuadTree[float64]

	players map[entities.PlayerID]*entities.Player

	bulletNextID uint32
	bullets      map[uint32]*entities.Bullet

	asteroidNextID uint32
	asteroids      map[uint32]*entities.Asteroid

	inputsQueue chan game.Input
	// for now we have one event: receive new player id
	incomingEvents        chan game.Event
	lastPlayerSnapshots   map[entities.PlayerID]game.PlayerState
	lastBulletSnapshots   map[uint32]game.ProjectileState
	lastAsteroidSnapshots map[uint32]game.ProjectileState

	lastTickWhenSpawn uint32

	currentSnapshot      *game.Snapshot
	tick                 uint32
	lastFullSnapshotTick uint32
	network              network.NetworkServer
	buf                  []byte
}

func (w *World) Init() {
	// world size
	w.Width = 500
	w.Height = 700
	topUnvisible := 200
	w.playerSpawnPoint = gmath.Vec{X: float64(w.Width) / 2, Y: float64(w.Height)/2 + float64(topUnvisible)}
	w.asteroidSpawnArea = gmath.Vec{X: float64(w.Width), Y: 200}

	// inital quad tree is world size
	worldRect := quadtree.NewRect[float64](0, 0, float64(w.Width), float64(w.Height))
	w.qtree = quadtree.New[float64](worldRect, 10)

	// initialize entites
	w.players = make(map[entities.PlayerID]*entities.Player)
	w.bullets = make(map[uint32]*entities.Bullet)
	w.asteroids = make(map[uint32]*entities.Asteroid)

	w.lastPlayerSnapshots = make(map[entities.PlayerID]game.PlayerState)
	w.lastBulletSnapshots = map[uint32]game.ProjectileState{}
	w.lastAsteroidSnapshots = map[uint32]game.ProjectileState{}
	w.inputsQueue = make(chan game.Input, 1024)
	w.incomingEvents = make(chan game.Event, 1024)
	w.currentSnapshot = new(game.Snapshot)
	w.tick = 0
	network, err := network.NewServer(os.Getenv("S_LISTEN"), os.Getenv("S_PORT"))
	if err != nil {
		// can't create udp server
		panic(err)
	}
	w.network = network
	w.buf = make([]byte, 2048)
}
