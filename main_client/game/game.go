package game

import (
	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/game/entities"
	"github.com/TaRosh/online_mover/main_client/asteroid"
	"github.com/TaRosh/online_mover/main_client/bullet"
	view "github.com/TaRosh/online_mover/main_client/camera"
	"github.com/TaRosh/online_mover/main_client/player"
	"github.com/TaRosh/online_mover/main_client/world"
	"github.com/TaRosh/online_mover/network"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	// get from inital connect
	WorldWidth  int
	WorldHeight int
	world       *world.World
	camera      *view.Camera

	Width  int
	Height int

	opt *ebiten.DrawImageOptions

	Network network.NetworkClient

	lastServerTick           uint32
	lastSnapshotForReconcile *game.Snapshot
	interpolationDelay       uint32

	snapshotQueue chan game.Snapshot
	// length of snapshot history
	maxSnapshot    int
	snapshotBuffer []*game.Snapshot
	events         chan game.Event
	inputsHistory  []game.Input

	Tick uint32
	buf  []byte
	// separate becouse we want recocilate
	localPlayer *player.Player
	players     map[entities.PlayerID]*player.Player
	// debugPlayer *Player
	projectiles map[uint32]*bullet.Bullet
	asteroids   map[uint32]*asteroid.Asteroid

	// input
	buttons uint8
	input   game.Input
}
