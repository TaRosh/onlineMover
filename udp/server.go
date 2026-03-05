package udp

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/TaRosh/online_mover/mover"
)

type server struct {
	conn *net.UDPConn

	done bool
	buf  []byte
	state
	userAddr *net.UDPAddr
}

const TickRate = 1

func (s *server) receive() {
	for {
		n, userAddr, err := s.conn.ReadFromUDP(s.buf)
		if err != nil {
			log.Println("server:receive", err)
			continue
		}
		s.userAddr = userAddr
		packet := Packet{}
		err = packet.Decode(s.buf[:n])
		if err != nil {
			log.Println("server:receive:", err)
		}
		s.processPacket(&packet)
	}
}

func (s *server) Listen() {
	defer s.Close()
	go s.receive()
	ticker := time.NewTicker(time.Second / TickRate)
	for range ticker.C {
		if s.userAddr == nil {
			continue
		}
		outPacket := NewPacket(s.state.id, s.state.lastIDReceived, s.state.packetsIGot, nil)
		data, err := outPacket.Encode()
		if err != nil {
			log.Println("Invalid data to send: ", err)
		}
		// echo back

		// fmt.Printf("Ack=%d AckBits=%032b\n", outPacket.Ack, outPacket.AckBits)
		fmt.Printf("[RECV] receved from client=%d AckBits:%b\n", outPacket.Ack, outPacket.AckBits)
		_, err = s.conn.WriteToUDP(data, s.userAddr)
		s.id += 1
		if err != nil {
			log.Println("server:Listen:write:", err)
		}
	}
}

func (s *server) processPacket(p *Packet) {
	if isNewer(p.Sequence, s.lastIDReceived) {
		shift := p.Sequence - s.lastIDReceived
		if shift >= 32 {
			s.packetsIGot = 0
		} else {
			s.packetsIGot <<= shift
		}
		s.packetsIGot |= 1
		s.lastIDReceived = p.Sequence
	} else {
		diff := s.lastIDReceived - p.Sequence
		if diff < 32 && diff >= 1 {
			bitIndex := diff - 1
			s.packetsIGot |= (1 << bitIndex)
		}
		fmt.Printf("Already got package: %d\n", p.Sequence)
		// if diff >= 32 {
		// 	// TODO: package to old what to do?
		// }
	}
	move := mover.Move{}
	move.Decode(p.Data)
	fmt.Printf("%+v\n", move)
}

func (s *server) Close() {
	s.conn.Close()
}

func NewServer(port string) (*server, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("", port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &server{
		conn: conn,
		state: state{
			id: 1,
		},
		buf: make([]byte, 2048),
	}, nil
}
