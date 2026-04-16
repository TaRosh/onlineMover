package game

import (
	"math"

	"github.com/quasilyte/gmath"
)

const twoPi = math.Pi * 2

func RadToUint16(r gmath.Rad) uint16 {
	r = r.Normalized()
	return uint16(r / twoPi * math.MaxUint16)
}

func Uint16ToRad(v uint16) float64 {
	return float64(v) / math.MaxUint16 * twoPi
}

