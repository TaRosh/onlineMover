package udp

import (
	"encoding/binary"
	"errors"
	"net"
	"time"
)

const headerSize = 13

type packetType uint8

const (
	InputPacket packetType = iota
	SnapshotPacket
)

type SentPacket struct {
	sendedWhen time.Time
	delivered  bool
}

type Header struct {
	Sequence uint32
	Ack      uint32
	AckBits  uint32
	Type     packetType
}

type Packet struct {
	Header
	Data []byte
	// addr we need when received package
	Addr *net.UDPAddr
}

func NewPacket(seq, ack, ackBits uint32, t packetType, data []byte) *Packet {
	header := Header{
		Sequence: seq,
		Ack:      ack,
		AckBits:  ackBits,
		Type:     t,
	}
	p := Packet{
		Header: header,
		Data:   data,
	}
	return &p
}

func (p *Packet) Encode(buf []byte) (int, error) {
	packetLen := headerSize + len(p.Data)
	if packetLen > len(buf) {
		return 0, errors.New("given buffer is smaller than packet size ( header + data )")
	}

	binary.BigEndian.PutUint32(buf[0:4], p.Header.Sequence)
	binary.BigEndian.PutUint32(buf[4:8], p.Header.Ack)
	binary.BigEndian.PutUint32(buf[8:12], p.Header.AckBits)
	buf[12] = byte(p.Header.Type)
	copy(buf[headerSize:], p.Data)
	return packetLen, nil
}

func (p *Packet) Decode(data []byte) error {
	header := data[:headerSize]
	if len(header) < headerSize {
		return errors.New("invalid packet: header size smaller then 12")
	}
	p.Header.Sequence = binary.BigEndian.Uint32(data[0:4])
	p.Header.Ack = binary.BigEndian.Uint32(data[4:8])
	p.Header.AckBits = binary.BigEndian.Uint32(data[8:12])
	p.Header.Type = packetType(data[12])

	// change buf allocation to get data completly
	// p.Data = make([]byte, len(data)-headerSize)
	dataSize := len(data) - headerSize
	if cap(p.Data) < dataSize {
		p.Data = make([]byte, dataSize)
	} else {
		p.Data = p.Data[:dataSize]
	}

	copy(p.Data, data[headerSize:])
	return nil
}

/*
Step 2

Add Sequence field only.

Detect packet loss.
*/
