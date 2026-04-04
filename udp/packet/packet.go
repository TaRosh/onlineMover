package packet

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type SentPacket struct {
	SendedWhen time.Time
	Delivered  bool
}

type Packet struct {
	Header
	Data []byte
	// addr we need when received package
	Addr *net.UDPAddr
}

func (p *Packet) Encode(buf []byte) (int, error) {
	packetLen := HeaderSize + len(p.Data)
	if packetLen > len(buf) {
		return 0, errors.New("given buffer is smaller than packet size ( header + data )")
	}
	n, _ := p.Header.Encode(buf)
	if HeaderSize != n {
		return 0, errors.New("encoded header size not match headerSize const")
	}

	copy(buf[n:], p.Data)
	return packetLen, nil
}

func (p *Packet) DecodeBody(data []byte) {
	if len(data) == 0 {
		p.Data = nil
	} else {
		p.Data = make([]byte, len(data))
		copy(p.Data, data)
	}
	return
}

func (p *Packet) Decode(data []byte) error {
	header := data[:HeaderSize]
	if len(header) < HeaderSize {
		return fmt.Errorf("Packet:Decode: %w", ErrHeaderSizeMismatch)
	}
	n, err := p.Header.Decode(data)
	if err != nil {
		return fmt.Errorf("Packet:Decode: %w", err)
	}

	p.DecodeBody(data[n:])

	return nil
}
