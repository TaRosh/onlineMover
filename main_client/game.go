package main

import (
	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/network"
)

type Game struct {
	Width  int
	Height int

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
	localPlayer *Player
	players     map[game.PlayerID]*Player
	// debugPlayer *Player
	projectiles map[uint32]*Bullet

	// input
	buttons uint8
	input   game.Input
}
