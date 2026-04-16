package network

import "encoding/binary"

type PlayerConnect struct {
	ID          uint32
	WorldWidth  uint16
	WorldHeight uint16
}

func (p *PlayerConnect) Encode(buf []byte) (int, error) {
	var offset int
	binary.BigEndian.PutUint32(buf[offset:], p.ID)
	offset += 4
	binary.BigEndian.PutUint16(buf[offset:], p.WorldWidth)
	offset += 2
	binary.BigEndian.PutUint16(buf[offset:], p.WorldHeight)
	offset += 2

	return offset, nil
}

func (p *PlayerConnect) Decode(data []byte) error {
	var offset int
	p.ID = binary.BigEndian.Uint32(data[offset:])
	offset += 4
	p.WorldWidth = binary.BigEndian.Uint16(data[offset:])
	offset += 2
	p.WorldHeight = binary.BigEndian.Uint16(data[offset:])
	offset += 2
	return nil
}
