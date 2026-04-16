package quadtree

type Bounder[T numeric] interface {
	Id() uint32
	Bounds() Rect[T]
	Layer() Layer
	Mask() Layer
}
