package quadtree

type QuadTree[N numeric] struct {
	boundary  Rect[N]
	capacity  int
	elems     []Bounder[N]
	divided   bool
	northwest *QuadTree[N]
	northeast *QuadTree[N]
	southwest *QuadTree[N]
	southeast *QuadTree[N]
}

func New[N numeric](boundary Rect[N], capacity int) *QuadTree[N] {
	return &QuadTree[N]{
		boundary: boundary,
		capacity: capacity,
		elems:    make([]Bounder[N], 0, capacity),
	}
}

func (qt *QuadTree[A]) Reset() {
	qt.elems = qt.elems[:0]
	if qt.divided {
		qt.northeast.Reset()
		qt.northwest.Reset()
		qt.southeast.Reset()
		qt.southwest.Reset()
	}
	qt.divided = false
	qt.northwest = nil
	qt.northeast = nil
	qt.northwest = nil
	qt.northwest = nil
}

func (qt *QuadTree[N]) Insert(elem Bounder[N]) bool {
	if !elem.Bounds().In(qt.boundary) {
		return false
	}
	if len(qt.elems) < qt.capacity {
		qt.elems = append(qt.elems, elem)
		return true
	} else {
		if !qt.divided {
			qt.subdivide()
		}
	}

	if qt.northeast.Insert(elem) {
		return true
	} else if qt.northwest.Insert(elem) {
		return true
	} else if qt.southeast.Insert(elem) {
		return true
	} else if qt.southwest.Insert(elem) {
		return true
	}
	return false
}

func (qt *QuadTree[N]) Query(screen Bounder[N]) []Bounder[N] {
	if !screen.Bounds().Overlaps(qt.boundary) {
		return nil
	}
	var found []Bounder[N]
	for _, elem := range qt.elems {
		// if elem.Bounds().Overlaps(screen) {
		if screen.Bounds().Overlaps(elem.Bounds()) {
			found = append(found, elem)
		}
	}
	if qt.divided {
		found = append(found, qt.northwest.Query(screen)...)
		found = append(found, qt.northeast.Query(screen)...)
		found = append(found, qt.southwest.Query(screen)...)
		found = append(found, qt.southeast.Query(screen)...)
	}
	return found
}

func (qt *QuadTree[N]) subdivide() {
	width := qt.boundary.Dx()
	height := qt.boundary.Dy()
	start := qt.boundary.Min
	nw := NewRect(start.X, start.Y, start.X+width/2, start.Y+height/2)
	ne := NewRect(start.X+width/2, start.Y, start.X+width, start.Y+height/2)
	sw := NewRect(start.X, start.Y+height/2, start.X+width/2, start.Y+height)
	se := NewRect(start.X+width/2, start.Y+height/2, start.X+width, start.Y+height)

	qt.northwest = New[N](nw, qt.capacity)
	qt.northeast = New[N](ne, qt.capacity)
	qt.southwest = New[N](sw, qt.capacity)
	qt.southeast = New[N](se, qt.capacity)
	qt.divided = true
}
