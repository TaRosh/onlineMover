package main

func stateSpawn(w *World) stateFn {
	if w.tick-w.lastTickWhenSpawn > 50 && len(w.asteroids) < 10 {
		w.spawnAsteroid(60)
		w.lastTickWhenSpawn = w.tick
	}
	return nil
}
