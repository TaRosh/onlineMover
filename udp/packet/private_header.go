package packet

import (
	"encoding/binary"
)

const PrivateHeaderSize = 13

type PrivateHeader struct {
	Ack     uint64
	AckBits uint32
	Type    Type
}

func (privH *PrivateHeader) Decode(data []byte) (int, error) {
	if len(data) < PrivateHeaderSize {
		return 0, ErrSmallBufferSize
	}
	var offset int
	privH.Ack = binary.BigEndian.Uint64(data[offset:])
	offset += 8
	privH.AckBits = binary.BigEndian.Uint32(data[offset:])
	offset += 4
	privH.Type = Type(data[12])
	offset += 1
	if offset != PrivateHeaderSize {
		return 0, ErrHeaderSizeMismatch
	}
	return offset, nil
}

func (privH *PrivateHeader) Encode(buf []byte) (int, error) {
	if len(buf) < PrivateHeaderSize {
		return 0, ErrSmallBufferSize
	}
	var offset int
	binary.BigEndian.PutUint64(buf[offset:], privH.Ack)
	offset += 8
	binary.BigEndian.PutUint32(buf[offset:], privH.AckBits)
	offset += 4
	buf[12] = byte(privH.Type)
	offset += 1
	return offset, nil
}
