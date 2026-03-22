package game

import (
	"encoding/binary"
)

const (
	InputLeft = 1 << iota
	InputRight
	InputUp
	InputDown
	InputShoot
)

// Input size total is 5 bytes
type Input struct {
	ID      PlayerID
	Tick    uint32
	Buttons uint8
}

func (i *Input) Encode(buf []byte) (int, error) {
	var offset int
	binary.BigEndian.PutUint32(buf[offset:], uint32(i.ID))
	offset += 4
	binary.BigEndian.PutUint32(buf[offset:], i.Tick)
	offset += 4
	buf[offset] = i.Buttons
	offset += 1

	return offset, nil
}

func (i *Input) Decode(data []byte) error {
	var offset int
	i.ID = PlayerID(binary.BigEndian.Uint32(data[offset:]))
	offset += 4
	i.Tick = binary.BigEndian.Uint32(data[offset:])
	offset += 4
	i.Buttons = data[offset]
	return nil
}
