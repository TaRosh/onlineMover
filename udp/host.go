package udp

import (
	"fmt"
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
	mutex       sync.Mutex
}
type host struct {
	conn *net.UDPConn

	receiveBuf []byte
	sentBuf    []byte

	// sentPacket    *Packet
	// receivePacket *Packet
	connections map[string]*connection

	mu sync.RWMutex
}

type Host interface {
	Sent(addr *net.UDPAddr, t packetType, data []byte) error
	Receive(sendPacketHere chan<- Packet) error
	GetAddr() *net.UDPAddr
	Close()
}

func (h *host) GetAddr() *net.UDPAddr {
	return h.conn.LocalAddr().(*net.UDPAddr)
}

func (h *host) Sent(addr *net.UDPAddr, t packetType, data []byte) error {
	packet := Packet{}

	conn := h.getConnection(addr)

	packet.Header.Sequence = conn.id
	packet.Header.Ack = conn.lastIDReceived
	packet.Header.AckBits = conn.packetsIGot
	packet.Type = t
	packet.Data = data
	n, err := packet.Encode(h.sentBuf)
	if err != nil {
		return err
	}
	_, err = h.conn.WriteToUDP(h.sentBuf[:n], addr)
	if err != nil {
		return err
	}
	conn.mutex.Lock()
	conn.packetsSend[conn.id] = SentPacket{
		sendedWhen: time.Now(),
		delivered:  false,
	}
	conn.id += 1
	conn.mutex.Unlock()
	return nil
}

func (c *connection) getSentPacket(id uint32) (SentPacket, bool) {
	p, exist := c.packetsSend[id]
	return p, exist
}

func (h *host) getConnection(addr *net.UDPAddr) *connection {
	h.mu.RLock()
	conn, exist := h.connections[addr.String()]
	h.mu.RUnlock()
	if !exist {

		c := connection{}
		c.id = 1
		c.packetsSend = make(map[uint32]SentPacket)
		h.mu.Lock()
		h.connections[addr.String()] = &c
		h.mu.Unlock()
		conn = &c
	}
	return conn
}

func (h *host) Receive(sendPacketHere chan<- Packet) error {
	n, userAddr, err := h.conn.ReadFromUDP(h.receiveBuf)
	if err != nil {
		return fmt.Errorf("transport:receive:%w", err)
	}

	conn := h.getConnection(userAddr)

	packet := Packet{}
	packet.Addr = userAddr
	err = packet.Decode(h.receiveBuf[:n])
	if err != nil {
		return fmt.Errorf("transport:receive:decode: %w", err)
	}
	conn.mutex.Lock()
	isDuplicate := h.processPacket(conn, &packet)
	conn.mutex.Unlock()
	// if packet double just drop
	// duplicate
	// continue
	if !isDuplicate {
		sendPacketHere <- packet
	}
	return nil
}

// check is packet id is old or new
func isNewer(a, b uint32) bool {
	return int32(a-b) > 0
}

func (c *connection) calcRtt(sample time.Duration) {
	if c.smoothedRTT == 0 {
		c.smoothedRTT = sample
	} else {
		c.smoothedRTT = time.Duration(float64(c.smoothedRTT)*0.9 + float64(sample)*0.1)
	}
}

func (c *connection) deleteSentPacket(id uint32) {
	delete(c.packetsSend, id)
}

func (h *host) processAck(conn *connection, header Header) {
	// check ack & ack bits for get rtt and calculate average rtt

	if sentPacket, ok := conn.getSentPacket(header.Ack); ok {
		if !sentPacket.delivered {
			rawRtt := time.Since(sentPacket.sendedWhen)
			conn.calcRtt(rawRtt)
			conn.deleteSentPacket(header.Ack)
		}
	}
	for i := range uint32(32) {
		if (header.AckBits>>i)&1 == 1 {
			if sentPacket, ok := conn.getSentPacket(header.Ack - i); ok {
				if !sentPacket.delivered {
					rawRtt := time.Since(sentPacket.sendedWhen)
					conn.calcRtt(rawRtt)
					conn.deleteSentPacket(header.Ack - i)
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
				conn.deleteSentPacket(id)
			}
		}
		conn.lastCleanup = time.Now()
	}
}

// return true if it's duplicate packet
func (h *host) processPacket(conn *connection, p *Packet) bool {
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
		if diff < 32 && diff >= 0 {
			bitIndex := diff
			// get double
			if conn.packetsIGot>>bitIndex&1 != 0 {
				return true
			}
			conn.packetsIGot |= (1 << bitIndex)
		}
		// if diff >= 32 {
		// 	// TODO: package to old what to do?
		// }
	}
	return false
}

func (h *host) Close() {
	h.conn.Close()
}

func NewClient(h, port string) (*host, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(h, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
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

func NewServer(h, port string) (*host, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(h, port))
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
