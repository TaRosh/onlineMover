package udp

import (
	"log"
	"net"
)

type server struct {
	conn *net.UDPConn

	done          bool
	receiveBuf    []byte
	sentBuf       []byte
	sentPacket    *Packet
	receivePacket *Packet
	state
	userAddr *net.UDPAddr
}

const TickRate = 1

type NetworkServer interface {
	SendSnapshot(data []byte) error
}

func (s *server) SendSnapshot(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	err := s.Write(SnapshotPacket, data)
	return err
}

func (s *server) Write(t packetType, data []byte) error {
	if s.userAddr == nil {
		return nil
	}
	if s.sentPacket == nil {
		s.sentPacket = new(Packet)
	}
	s.sentPacket.Header.Sequence = s.state.id
	s.sentPacket.Header.Ack = s.state.lastIDReceived
	s.sentPacket.Header.AckBits = s.state.packetsIGot
	s.sentPacket.Type = t
	s.sentPacket.Data = data
	n, err := s.sentPacket.Encode(s.sentBuf)
	if err != nil {
		return err
	}
	_, err = s.conn.WriteToUDP(s.sentBuf[:n], s.userAddr)
	if err != nil {
		return err
	}
	s.state.id += 1
	return nil
}

func (s *server) Receive(sendPacketHere chan<- Packet) {
	if s.receivePacket == nil {
		s.receivePacket = new(Packet)
	}
	for {
		n, userAddr, err := s.conn.ReadFromUDP(s.receiveBuf)
		if err != nil {
			log.Println("server:receive", err)
			continue
		}
		s.userAddr = userAddr
		err = s.receivePacket.Decode(s.receiveBuf[:n])
		if err != nil {
			log.Println("server:receive:", err)
		}
		s.processPacket(s.receivePacket)
		sendPacketHere <- *s.receivePacket
	}
}

// incomingPackets := make(<-chan udp.Packet, 5)

func (s *server) processPacket(p *Packet) {
	// if packet id is new ( < then last id i got)
	// check is inside our check table ( uint32 ) or we lost
	// 32 packet and it is newer
	// do new check table or shift for it's diff ( currentId - lastId )
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
		// if diff >= 32 {
		// 	// TODO: package to old what to do?
		// }
	}
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
		receiveBuf: make([]byte, 2048),
		sentBuf:    make([]byte, 2048),
	}, nil
}
