package game

import (
	"os"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/game/entities"
	"github.com/TaRosh/online_mover/main_client/asteroid"
	"github.com/TaRosh/online_mover/main_client/bullet"
	"github.com/TaRosh/online_mover/main_client/camera"
	"github.com/TaRosh/online_mover/main_client/player"
	"github.com/TaRosh/online_mover/main_client/world"
	"github.com/TaRosh/online_mover/network"
	"github.com/quasilyte/gmath"
)

func (g *Game) Init() {
	g.Width, g.Height = 500, 500

	g.world = world.New()

	g.camera = camera.New(float64(g.Width), float64(g.Height))
	g.camera.Pos = gmath.Vec{X: 0, Y: 200}

	g.players = make(map[entities.PlayerID]*player.Player)
	g.projectiles = make(map[uint32]*bullet.Bullet)
	g.asteroids = make(map[uint32]*asteroid.Asteroid)

	network, err := network.NewClient(os.Getenv("S_HOST"), os.Getenv("S_PORT"))
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
