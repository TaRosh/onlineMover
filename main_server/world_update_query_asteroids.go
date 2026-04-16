package main

import (
	"fmt"

	"github.com/TaRosh/online_mover/game/layers"
)

// find entities for collision check
func (w *World) queryAsteroids() {
	for _, a := range w.asteroids {
		near := w.qtree.Query(a)
		for _, other := range near {
			if other.Id() == a.Id() && a.Layer() == other.Layer() {
				continue
			}
			// check layer
			// mask give what we interact with
			// layer give entite type
			if a.Mask()&other.Layer() == 0 {
				continue
			}

			if a.Bounds().Overlaps(other.Bounds()) {
				// TODO: what do on collision
				// asteroid collide with player or bullet
				switch other.Layer() {
				case layers.LayerPlayer:
					fmt.Println("queryAsteroids:HIT")
					w.asteroids[a.Id()].Deleted = true

				case layers.LayerBullet:
					fmt.Println("queryAsteroids:BULLET HIT")
					w.asteroids[a.Id()].Deleted = true
					w.bullets[other.Id()].Deleted = true
				}
			}
		}
	}
}
