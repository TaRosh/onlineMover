package game

import (
	"encoding/binary"
	"math"

	"github.com/quasilyte/gmath"
)

// size 4 + 8 * 2
type Input struct {
	Tick uint32
	Vel  gmath.Vec
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func (i *Input) Encode(buf []byte) (int, error) {
	binary.BigEndian.PutUint32(buf[0:4], i.Tick)
	binary.BigEndian.PutUint64(buf[4:12], math.Float64bits(i.Vel.X))
	binary.BigEndian.PutUint64(buf[12:20], math.Float64bits(i.Vel.Y))

	return 20, nil
}

func (i *Input) Decode(data []byte) error {
	i.Tick = binary.BigEndian.Uint32(data[0:4])
	i.Vel.X = math.Float64frombits(binary.BigEndian.Uint64(data[4:12]))
	i.Vel.Y = math.Float64frombits(binary.BigEndian.Uint64(data[12:20]))
	return nil
}
