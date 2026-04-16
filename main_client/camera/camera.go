package camera

import (
	"github.com/quasilyte/gmath"
)

type Camera struct {
	Pos    gmath.Vec
	Width  float64
	Height float64
}

func New(w, h float64) *Camera {
	c := Camera{
		Width:  w,
		Height: h,
	}
	return &c
}
