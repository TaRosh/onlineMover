package udp

import (
	"encoding/binary"
	"errors"
	"net"
	"time"
)

type SentPacket struct {
	sendedWhen time.Time
	delivered  bool
}

type Header struct {
	Sequence uint32
	Ack      uint32
	AckBits  uint32
}

type Packet struct {
	Header
	Data []byte
	// addr we need when received package
	Addr *net.UDPAddr
}

func NewPacket(seq, ack, ackBits uint32, data []byte) *Packet {
	header := Header{
		Sequence: seq,
		Ack:      ack,
		AckBits:  ackBits,
	}
	p := Packet{
		Header: header,
		Data:   data,
	}
	return &p
}

func (p *Packet) Encode() ([]byte, error) {
	buf := make([]byte, headerSize+len(p.Data))

	binary.BigEndian.PutUint32(buf[0:4], p.Header.Sequence)
	binary.BigEndian.PutUint32(buf[4:8], p.Header.Ack)
	binary.BigEndian.PutUint32(buf[8:12], p.Header.AckBits)
	copy(buf[headerSize:], p.Data)
	return buf, nil
}

func (p *Packet) Decode(data []byte) error {
	header := data[:headerSize]
	if len(header) < headerSize {
		return errors.New("invalid packet: header size smaller then 12")
	}
	p.Header.Sequence = binary.BigEndian.Uint32(data[0:4])
	p.Header.Ack = binary.BigEndian.Uint32(data[4:8])
	p.Header.AckBits = binary.BigEndian.Uint32(data[8:12])
	p.Data = make([]byte, len(data)-headerSize)
	copy(p.Data, data[headerSize:])
	return nil
}

/*
Step 2

Add Sequence field only.

Detect packet loss.
*/
