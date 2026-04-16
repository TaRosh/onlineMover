package main

func (w *World) buildQt() {
	w.qtree.Reset()
	for _, p := range w.players {
		w.qtree.Insert(p)
	}
	for _, b := range w.bullets {
		if b.Deleted {
			delete(w.bullets, b.Id())
			continue
		}
		w.qtree.Insert(b)
	}
	for _, a := range w.asteroids {
		if a.Deleted {
			delete(w.asteroids, a.Id())
			continue
		}
		w.qtree.Insert(a)
	}
}
