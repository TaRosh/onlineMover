package main

import (
	"math/rand/v2"

	"github.com/TaRosh/online_mover/game/entities"
	"github.com/quasilyte/gmath"
)

func (w *World) spawnAsteroid(radius float64) {
	if len(w.players) <= 0 {
		return
	}

	a := entities.NewAsteroid(w.asteroidNextID, radius, radius)
	spawnX := rand.Float64() * w.asteroidSpawnArea.X
	if spawnX == 0 {
		spawnX += radius
	} else if spawnX+radius > w.asteroidSpawnArea.X {
		spawnX -= radius
	}
	spawnY := rand.Float64() * w.asteroidSpawnArea.Y
	if spawnY+radius < 0 {
		spawnY += 0 + radius
	} else if spawnY+radius > w.asteroidSpawnArea.Y {
		spawnY = w.asteroidSpawnArea.Y - radius
	}
	a.Position = gmath.Vec{
		X: spawnX,
		Y: spawnY,
	}
	// steer = (desired - vel ) * maxForce
	var steer gmath.Vec
	for _, p := range w.players {
		des := a.Position.DirectionTo(p.Position).Mulf(a.MaxForce)
		des = des.Sub(a.Velocity)
		steer = steer.Add(des)
	}
	a.ApplyForce(steer)
	w.asteroids[w.asteroidNextID] = a
	w.asteroidNextID += 1
}
