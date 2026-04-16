package game

import (
	"image/color"
	"log"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/main_client/asteroid"
	"github.com/TaRosh/online_mover/main_client/bullet"
	"github.com/TaRosh/online_mover/main_client/converter"
	inner "github.com/TaRosh/online_mover/main_client/player"
)

func (g *Game) HandleDeltaSnapshot(snapshot *game.Snapshot) {
	// create new player ( if exist )
	// & reset player pos ( state )
	var err error
	for _, p := range snapshot.Players {
		player, exist := g.players[p.ID]
		if !exist {
			player, err = inner.New(p.ID, color.RGBA{0xff, 0, 0, 0xff})
			if err != nil {
				log.Println("handelFullSnapshot:NewPlayer:", err)
			}
			g.players[p.ID] = player
		}
		converter.StateToPlayer(player, p)
	}

	// create bullets
	for _, projectile := range snapshot.Projectiles {
		b, exist := g.projectiles[projectile.ID]
		if !exist {
			b = bullet.NewBullet(projectile.ID, color.RGBA{0, 0, 0xff, 0xff})
			g.projectiles[b.ID] = b
		}
		converter.StateToBullet(b, projectile)
	}
	// create asteroids
	for _, ast := range snapshot.Asteroids {
		a, exist := g.asteroids[ast.ID]
		if !exist {
			a, err = asteroid.New(ast.ID, 60, color.White)
			if err != nil {
				// TODO: whate happen?
				panic(err)
			}
			g.asteroids[a.ID] = a
		}
		converter.StateToAsteroid(a, ast)
		// time.Sleep(time.Second)
	}
}
