package quadtree

type point[T numeric] struct {
	X T
	Y T
}

type Rect[T numeric] struct {
	Min point[T]
	Max point[T]
}

func (r Rect[T]) Bounds() Rect[T] {
	return r
}

func (r Rect[T]) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

func (r Rect[T]) In(s Rect[T]) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	// return s.Max.X <= r.Min.X &&
	// 	s.Max.Y <= r.Min.Y
	return s.Min.X <= r.Min.X && r.Min.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Min.Y <= s.Max.Y
}

func (r Rect[T]) Overlaps(s Rect[T]) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

func (r Rect[T]) Dx() T {
	return r.Max.X - r.Min.X
}

func (r Rect[T]) Dy() T {
	return r.Max.Y - r.Min.Y
}

func NewRect[N numeric](minX, minY, maxX, maxY N) Rect[N] {
	return Rect[N]{
		Min: point[N]{
			X: minX,
			Y: minY,
		},
		Max: point[N]{
			X: maxX,
			Y: maxY,
		},
	}
}
