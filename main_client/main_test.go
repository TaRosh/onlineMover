package main

import (
	"testing"

	"github.com/TaRosh/online_mover/game"
	"github.com/quasilyte/gmath"
	"github.com/stretchr/testify/require"
)

func TestInterpolationFromSnapshot(t *testing.T) {
	/*
		tick100 → pos(0,0)
		tick102 → pos(10,0)
		renderTick = 101
	*/
	// Test: simple interpolation: alpha = 0.5
	g := Game{
		lastServerTick:     102,
		interpolationDelay: 1,
	}
	playerID := game.PlayerID(1)
	g.snapshotBuffer = []*game.Snapshot{
		{
			Tick: 100,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  0,
					Y:  0,
				},
			},
		},
		{
			Tick: 102,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  1000,
					Y:  0,
				},
			},
		},
	}
	player := &Player{
		Player: game.Player{
			ID: playerID,
		},
	}
	pos := g.playerInterpolation(g.snapshotBuffer[0], g.snapshotBuffer[1], 0.5, player)
	require.Equal(t, float64(5), pos.X)

	/*
		tick100 → (0,0)
		tick101 → (10,0)
		renderTick = 100
	*/
	// Test: alpha 0 ( start snapshot )
	g.lastServerTick = 101
	g.interpolationDelay = 1
	g.snapshotBuffer = []*game.Snapshot{
		{
			Tick: 100,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  0,
					Y:  0,
				},
			},
		},
		{
			Tick: 101,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  1000,
					Y:  0,
				},
			},
		},
	}
	pos = g.playerInterpolation(g.snapshotBuffer[0], g.snapshotBuffer[1], 0, player)
	require.Equal(t, float64(0), pos.X)

	// Test: alpha 1: end snapshot
	g.lastServerTick = 102
	g.interpolationDelay = 1
	g.snapshotBuffer = []*game.Snapshot{
		{
			Tick: 100,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  0,
					Y:  0,
				},
			},
		},
		{
			Tick: 101,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  1000,
					Y:  0,
				},
			},
		},
	}
	pos = g.playerInterpolation(g.snapshotBuffer[0], g.snapshotBuffer[1], 1, player)
	require.Equal(t, float64(10), pos.X)
	// Test: diagonal move
	g.lastServerTick = 102
	g.interpolationDelay = 1
	g.snapshotBuffer = []*game.Snapshot{
		{
			Tick: 100,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  0,
					Y:  0,
				},
			},
		},
		{
			Tick: 102,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  1000,
					Y:  1000,
				},
			},
		},
	}
	pos = g.playerInterpolation(g.snapshotBuffer[0], g.snapshotBuffer[1], 0.5, player)
	require.Equal(t, float64(5), pos.X)
	require.Equal(t, float64(5), pos.Y)

	// Test: no snapshot case
	g.snapshotBuffer = make([]*game.Snapshot, 0)
	pos = g.playerInterpolation(nil, nil, 1, player)
	require.Equal(t, pos, player.Position)

	// Test: no found second snapshot
	g.lastServerTick = 200
	g.interpolationDelay = 1
	g.snapshotBuffer = []*game.Snapshot{
		{
			Tick: 100,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  0,
					Y:  0,
				},
			},
		},
		{
			Tick: 102,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  1000,
					Y:  1000,
				},
			},
		},
	}
	player.Position = gmath.Vec{X: 10, Y: 10}
	pos = g.playerInterpolation(g.snapshotBuffer[0], g.snapshotBuffer[1], 1, player)
	require.Equal(t, pos, player.Position)

	// Test: fractional alpha test
	/*
		tick100 → (0,0)
		tick104 → (8,0)
		renderTick = 101
		----------
		(101-100)/(104-100) = 0.25
	*/
	g.lastServerTick = 102
	g.interpolationDelay = 1
	g.snapshotBuffer = []*game.Snapshot{
		{
			Tick: 100,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  0,
					Y:  0,
				},
			},
		},
		{
			Tick: 104,
			Players: []game.PlayerState{
				{
					ID: playerID,
					X:  800,
					Y:  0,
				},
			},
		},
	}
	pos = g.playerInterpolation(g.snapshotBuffer[0], g.snapshotBuffer[1], 0.25, player)
	require.Equal(t, float64(2), pos.X)

	// TODO: Test: snapshot not contain player
}
