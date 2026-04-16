package main

func (w *World) update() {
	// build qtree
	w.buildQt()
	// TODO
	// 1. collision detect
	w.collisionDetect()
	// 2. collision resolve
	// 3. collision delete

	// query entities

	// update all
	for _, p := range w.players {
		p.Update()
		if p.Position.X-p.Width/2 < 0 {
			p.Position.X = 0 + p.Width/2
		} else if p.Position.X+p.Width/2 > float64(w.Width) {
			p.Position.X = float64(w.Width) - p.Width/2
		}
		if p.Position.Y-p.Height/2 < 0 {
			p.Position.Y = 0 + p.Height/2
		}
		if p.Position.Y+p.Height/2 > float64(w.Height) {
			p.Position.Y = float64(w.Height) - p.Height/2
		}
	}
	for _, b := range w.bullets {
		b.Update()
		if b.Position.X+b.Width < 0 {
			b.Deleted = true
		} else if b.Position.X-b.Width > float64(w.Width) {
			b.Deleted = true
		}
		if b.Position.Y+b.Height < 0 {
			b.Deleted = true
		}
		if b.Position.Y-b.Height > float64(w.Height) {
			b.Deleted = true
		}
	}
	for _, a := range w.asteroids {
		a.Update()
		if a.Position.X+a.Width < 0 {
			a.Deleted = true
		} else if a.Position.X-a.Width > float64(w.Width) {
			a.Deleted = true
		}
		if a.Position.Y+a.Height < 0 {
			a.Deleted = true
		}
		if a.Position.Y-a.Height > float64(w.Height) {
			a.Deleted = true
		}
	}
}
