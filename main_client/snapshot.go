package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/TaRosh/online_mover/game"
)

func (g *Game) handleFullSnapshot(snapshot *game.Snapshot) {
	seen := make(map[game.PlayerID]bool)
	// create new player ( if exist )
	// & reset player pos ( state )
	for _, p := range snapshot.Players {
		player, exist := g.players[p.ID]
		if !exist {
			player, err := NewPlayer(p.ID, color.RGBA{0xff, 0, 0, 0xff})
			if err != nil {
				log.Println("handelFullSnapshot:NewPlayer:", err)
			}
			g.players[p.ID] = player
		}
		stateToPlayer(player, p)
		seen[p.ID] = true
	}
	for id := range g.players {
		if !seen[id] {
			delete(g.players, id)
		}
	}

	// create bullets
	for _, bullet := range snapshot.Projectiles {
		b, exist := g.projectiles[bullet.ID]
		if !exist {
			b = NewBullet(bullet.ID, color.RGBA{0, 0, 0xff, 0xff})
			g.projectiles[b.ID] = b
		}
		stateToBullet(b, bullet)
	}
}

func (g *Game) handleDeltaSnapshot(snapshot *game.Snapshot) {
	// create new player ( if exist )
	// & reset player pos ( state )
	var err error
	for _, p := range snapshot.Players {
		player, exist := g.players[p.ID]
		if !exist {
			player, err = NewPlayer(p.ID, color.RGBA{0xff, 0, 0, 0xff})
			if err != nil {
				log.Println("handelFullSnapshot:NewPlayer:", err)
			}
			g.players[p.ID] = player
		}
		fmt.Println("TRY PROCESS PLAYER:", player)
		stateToPlayer(player, p)
	}

	// create bullets
	for _, bullet := range snapshot.Projectiles {
		b, exist := g.projectiles[bullet.ID]
		if !exist {
			b = NewBullet(bullet.ID, color.RGBA{0, 0, 0xff, 0xff})
			g.projectiles[b.ID] = b
		}
		stateToBullet(b, bullet)
	}
}
