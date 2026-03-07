package game

import (
	"encoding/binary"
)

const (
	InputLeft = 1 << iota
	InputRight
	InputUp
	InputDown
)

// Input size total is 5 bytes
type Input struct {
	Tick    uint32
	Buttons uint8
}

func (i *Input) Encode(buf []byte) (int, error) {
	binary.BigEndian.PutUint32(buf[0:4], i.Tick)
	buf[4] = i.Buttons

	return 5, nil
}

func (i *Input) Decode(data []byte) error {
	i.Tick = binary.BigEndian.Uint32(data[0:4])
	i.Buttons = data[4]
	return nil
}
