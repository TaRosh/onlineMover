package packet

import "encoding/binary"

const PublicHeaderSize = 12

type PublicHeader struct {
	ConnectionID uint32
	// TODO: implement quic hide full sequence
	// part of packet seq
	// SeqShort uint16
	Sequence uint64
}

func (pubH *PublicHeader) Encode(buf []byte) (int, error) {
	if len(buf) < PublicHeaderSize {
		return 0, ErrSmallBufferSize
	}
	var offset int
	binary.BigEndian.PutUint32(buf[offset:], pubH.ConnectionID)
	offset += 4
	binary.BigEndian.PutUint64(buf[offset:], pubH.Sequence)
	offset += 8
	return offset, nil
}

func (pubH *PublicHeader) Decode(data []byte) (int, error) {
	if len(data) < PublicHeaderSize {
		return 0, ErrSmallBufferSize
	}
	var offset int
	pubH.ConnectionID = binary.BigEndian.Uint32(data[offset:])
	offset += 4
	pubH.Sequence = binary.BigEndian.Uint64(data[offset:])
	offset += 8
	if offset != PublicHeaderSize {
		return 0, ErrHeaderSizeMismatch
	}
	return offset, nil
}
