package udp

import (
	"log"
	"net"
	"sync"
	"time"
)

type state struct {
	id             uint32
	lastIDReceived uint32
	packetsIGot    uint32
}

type connection struct {
	state
	Addr        *net.UDPAddr
	packetsSend map[uint32]SentPacket
	smoothedRTT time.Duration
	lastCleanup time.Time
}

type host struct {
	conn *net.UDPConn

	receiveBuf []byte
	sentBuf    []byte

	sentPacket    *Packet
	receivePacket *Packet
	connections   map[string]*connection

	mu sync.Mutex
}

type Host interface {
	Sent(addr *net.UDPAddr, t packetType, data []byte) error
	Receive(sendPacketHere chan<- Packet)
}

/*
 */
func (h *host) Sent(addr *net.UDPAddr, t packetType, data []byte) error {
	if h.sentPacket == nil {
		h.sentPacket = new(Packet)
	}
	conn, exist := h.connections[addr.String()]
	if !exist {
		//
		if !exist {
			// first time request to us
			// create connection abstraction
			conn = &connection{}
			conn.id = 1
			conn.packetsSend = make(map[uint32]SentPacket)
			h.connections[addr.String()] = conn
		}
	}
	h.sentPacket.Header.Sequence = conn.id
	h.sentPacket.Header.Ack = conn.lastIDReceived
	h.sentPacket.Header.AckBits = conn.packetsIGot
	h.sentPacket.Type = t
	h.sentPacket.Data = data
	n, err := h.sentPacket.Encode(h.sentBuf)
	if err != nil {
		return err
	}
	_, err = h.conn.WriteToUDP(h.sentBuf[:n], addr)
	if err != nil {
		return err
	}
	// TODO: Need mutex ?
	h.mu.Lock()
	conn.packetsSend[conn.id] = SentPacket{
		sendedWhen: time.Now(),
		delivered:  false,
	}
	h.mu.Unlock()
	conn.id += 1
	return nil
}

func (h *host) Receive(sendPacketHere chan<- Packet) {
	if h.receivePacket == nil {
		h.receivePacket = new(Packet)
	}
	for {
		n, userAddr, err := h.conn.ReadFromUDP(h.receiveBuf)
		if err != nil {
			log.Println("server:receive", err)
			continue
		}

		conn, exist := h.connections[userAddr.String()]
		if !exist {
			// first time request to us
			// create connection abstraction
			conn = &connection{}
			conn.id = 1
			conn.packetsSend = make(map[uint32]SentPacket)
			h.connections[userAddr.String()] = conn
		}
		h.receivePacket.Addr = userAddr
		err = h.receivePacket.Decode(h.receiveBuf[:n])
		if err != nil {
			log.Println("server:receive:", err)
		}
		h.processPacket(conn, h.receivePacket)
		// if packet doble just drop

		// continue
		sendPacketHere <- *h.receivePacket
	}
}

// check is packet id is old or new
func isNewer(a, b uint32) bool {
	return int32(a-b) > 0
}

func (h *host) calcRtt(conn *connection, sample time.Duration) {
	if conn.smoothedRTT == 0 {
		conn.smoothedRTT = sample
	} else {
		conn.smoothedRTT = time.Duration(float64(conn.smoothedRTT)*0.9 + float64(sample)*0.1)
	}
}

func (h *host) processAck(conn *connection, header Header) {
	// check ack & ack bits for get rtt and calculate average rtt
	if sentPacket, ok := conn.packetsSend[header.Ack]; ok {
		if !sentPacket.delivered {
			rawRtt := time.Since(sentPacket.sendedWhen)
			h.calcRtt(conn, rawRtt)
			delete(conn.packetsSend, header.Ack)
		}
	}
	for i := range uint32(32) {
		if (header.AckBits>>i)&1 == 1 {
			if sentPacket, ok := conn.packetsSend[header.Ack-i]; ok {
				if !sentPacket.delivered {
					rawRtt := time.Since(sentPacket.sendedWhen)
					h.calcRtt(conn, rawRtt)
					delete(conn.packetsSend, header.Ack)
				}
			}
		}
	}
	// clean up the map for sended packets
	// delete data about packet if  not arrived after 3 second
	// launch every second
	if time.Since(conn.lastCleanup) > time.Second {
		for id, sentPacket := range conn.packetsSend {
			if time.Since(sentPacket.sendedWhen) > time.Second*3 {
				delete(conn.packetsSend, id)
			}
		}
		conn.lastCleanup = time.Now()
	}
}

func (h *host) processPacket(conn *connection, p *Packet) {
	h.processAck(conn, p.Header)
	// if packet id is new ( < then last id i got)
	// check is inside our check table ( uint32 ) or we lost
	// 32 packet and it is newer
	// do new check table or shift for it's diff ( currentId - lastId )
	if isNewer(p.Sequence, conn.lastIDReceived) {
		shift := p.Sequence - conn.lastIDReceived
		if shift >= 32 {
			conn.packetsIGot = 0
		} else {
			conn.packetsIGot <<= shift
		}
		conn.packetsIGot |= 1
		conn.lastIDReceived = p.Sequence
	} else {
		diff := conn.lastIDReceived - p.Sequence
		if diff < 32 && diff >= 1 {
			bitIndex := diff
			conn.packetsIGot |= (1 << bitIndex)
		}
		// if diff >= 32 {
		// 	// TODO: package to old what to do?
		// }
	}
}

func (h *host) Close() {
	h.conn.Close()
}

func NewServer(port string) (*host, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("", port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &host{
		conn:        conn,
		connections: make(map[string]*connection),
		receiveBuf:  make([]byte, 2048),
		sentBuf:     make([]byte, 2048),
	}, nil
}
