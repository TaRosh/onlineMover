package udp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

const headerSize = 19

type packetType uint8

const (
	PacketInvalid packetType = iota
	PacketInput
	PacketSnapshot
	PacketConnect
	PacketKeyExchangeRequest
	PacketKeyExchangeAnswer
	PacketAccept
)

type SentPacket struct {
	sendedWhen time.Time
	delivered  bool
}

type PublicHeader struct {
	ConnectionID uint32
	// part of packet seq
	SeqShort uint16
}

type PrivateHeader struct {
	Sequence uint32
	Ack      uint32
	AckBits  uint32
	Type     packetType
}

type Header struct {
	PublicHeader
	PrivateHeader
}

type Packet struct {
	Header
	Data []byte
	// addr we need when received package
	Addr *net.UDPAddr
}

func (pubH *PublicHeader) Encode(buf []byte) int {
	var offset int
	binary.BigEndian.PutUint32(buf[offset:], pubH.ConnectionID)
	offset += 4
	binary.BigEndian.PutUint16(buf[offset:], pubH.SeqShort)
	offset += 2
	return offset
}

func (pubH *PublicHeader) Decode(data []byte) int {
	var offset int
	pubH.ConnectionID = binary.BigEndian.Uint32(data[offset:])
	offset += 4
	pubH.SeqShort = binary.BigEndian.Uint16(data[offset:])
	offset += 2
	return offset
}

func (privH *PrivateHeader) Decode(data []byte) int {
	var offset int
	privH.Sequence = binary.BigEndian.Uint32(data[offset:])
	offset += 4
	privH.Ack = binary.BigEndian.Uint32(data[offset:])
	offset += 4
	privH.AckBits = binary.BigEndian.Uint32(data[offset:])
	offset += 4
	privH.Type = packetType(data[12])
	offset += 1
	return offset
}

func (privH *PrivateHeader) Encode(buf []byte) int {
	var offset int
	binary.BigEndian.PutUint32(buf[offset:], privH.Sequence)
	offset += 4
	binary.BigEndian.PutUint32(buf[offset:], privH.Ack)
	offset += 4
	binary.BigEndian.PutUint32(buf[offset:], privH.AckBits)
	offset += 4
	buf[12] = byte(privH.Type)
	offset += 1
	return offset
}

func (h *Header) Encode(buf []byte) int {
	var totalOffset int
	n := h.PublicHeader.Encode(buf)
	totalOffset += n
	n = h.PrivateHeader.Encode(buf[n:])
	totalOffset += n
	return totalOffset
}

func (p *Packet) Encode(buf []byte) (int, error) {
	packetLen := headerSize + len(p.Data)
	if packetLen > len(buf) {
		return 0, errors.New("given buffer is smaller than packet size ( header + data )")
	}
	n := p.Header.Encode(buf)
	if headerSize != n {
		return 0, errors.New("encoded header size not match headerSize const")
	}

	copy(buf[n:], p.Data)
	return packetLen, nil
}

func (h *Header) Decode(data []byte) (int, error) {
	if len(data) < headerSize {
		return 0, errors.New("Header:Decode: buffer size smaller than expected header size")
	}
	var totalOffset int
	n := h.PublicHeader.Decode(data)
	totalOffset += n
	n = h.PrivateHeader.Decode(data[n:])
	totalOffset += n
	if totalOffset != headerSize {
		return 0, errors.New("Header:Decode: offset after decode don't match header size cost")
	}
	return totalOffset, nil
}

func (p *Packet) Decode(data []byte) error {
	header := data[:headerSize]
	if len(header) < headerSize {
		return errors.New("Packet:Decode: invalid packet: header size smaller then 12")
	}
	n, err := p.Header.Decode(data)
	if err != nil {
		return fmt.Errorf("Packet:Decode: %w", err)
	}

	dataSize := len(data) - headerSize
	if dataSize == 0 {
		p.Data = nil
	} else {
		p.Data = make([]byte, dataSize)
		copy(p.Data, data[n:])
	}

	return nil
}

/*
Step 2

Add Sequence field only.

Detect packet loss.
*/
