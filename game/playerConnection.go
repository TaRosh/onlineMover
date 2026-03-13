package game

import (
	"encoding/binary"
)

// PlayerState size is 36
type PlayerIDPacket struct {
	ID uint32
}

func (pId *PlayerIDPacket) Encode(buf []byte) (int, error) {
	var offset int
	binary.BigEndian.PutUint32(buf[0:], pId.ID)
	offset += 4
	return offset, nil
}

func (pId *PlayerIDPacket) Decode(data []byte) error {
	// var offset int
	pId.ID = binary.BigEndian.Uint32(data[0:])
	// offset += 4

	return nil
}
