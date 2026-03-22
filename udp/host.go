package udp

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type (
	ServerHost interface {
		Host
		Send(addr *net.UDPAddr, t packetType, data []byte) error
		CheckTimeouts(timeoutConnectionsHere chan<- *net.UDPAddr)
		DeleteConn(id *net.UDPAddr)
	}
	ClientHost interface {
		Host
		Send(t packetType, data []byte) error
	}
)

type Host interface {
	Receive(sendPacketHere chan<- Packet) error
	GetAddr() *net.UDPAddr
	Close()
}

const (
	bufferSize         = 1200
	serverKeyForClient = "server"
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

	// where we last time receive packet from that connection
	lastHeared  time.Time
	lastCleanup time.Time
	mutex       sync.RWMutex
}

type host struct {
	conn *net.UDPConn

	bufferPool sync.Pool
	packetPool sync.Pool

	// sentPacket    *Packet
	// receivePacket *Packet
	connections map[string]*connection

	mu sync.RWMutex
}

type serverHost struct {
	*host
}

type clientHost struct {
	*host
}

func (sh *serverHost) Send(addr *net.UDPAddr, t packetType, data []byte) error {
	key := addr.String()
	conn := sh.getConnection(key, addr)
	err := sh.send(conn, t, data)
	return err
}

func (ch *clientHost) Send(t packetType, data []byte) error {
	key := serverKeyForClient
	conn := ch.getConnection(key, nil)
	err := ch.send(conn, t, data)
	return err
}

func (sh *serverHost) Receive(sendPacketHere chan<- Packet) error {
	buf := sh.bufferPool.Get().([]byte)
	defer func() {
		buf = buf[:cap(buf)]
		sh.bufferPool.Put(buf)
	}()
	n, playerAddr, err := sh.conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	key := playerAddr.String()
	conn := sh.getConnection(key, playerAddr)
	fmt.Println("RECV packet")
	err = sh.processRawPacket(conn, buf[:n], sendPacketHere)
	if err != nil {
		return err
	}
	return nil
}

func (ch *clientHost) Receive(sendPacketHere chan<- Packet) error {
	buf := ch.bufferPool.Get().([]byte)
	defer func() {
		buf = buf[:cap(buf)]
		ch.bufferPool.Put(buf)
	}()
	n, serverAddr, err := ch.conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	key := serverKeyForClient
	conn := ch.getConnection(key, serverAddr)
	err = ch.processRawPacket(conn, buf[:n], sendPacketHere)
	if err != nil {
		return err
	}
	return nil
}

func (h *host) GetAddr() *net.UDPAddr {
	return h.conn.LocalAddr().(*net.UDPAddr)
}

func (h *host) DeleteConn(id *net.UDPAddr) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if id == nil {
		panic("delete conn")
	}
	delete(h.connections, id.String())
}

func (h *host) CheckTimeouts(noAnswerAddrHere chan<- *net.UDPAddr) {
	defer close(noAnswerAddrHere)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, conn := range h.connections {
		conn.mutex.RLock()
		lastTimeHeared := conn.lastHeared
		fmt.Println(time.Since(lastTimeHeared))
		conn.mutex.RUnlock()
		if time.Since(lastTimeHeared) > 5*time.Second {
			if conn.Addr == nil {
				panic("checkTimeouts")
			}
			noAnswerAddrHere <- conn.Addr
		}
	}
}

func (h *host) send(conn *connection, t packetType, data []byte) error {
	packet := h.packetPool.Get().(*Packet)
	defer func() {
		*packet = Packet{}
		h.packetPool.Put(packet)
	}()

	conn.mutex.RLock()
	packet.Header.Sequence = conn.id
	packet.Header.Ack = conn.lastIDReceived
	packet.Header.AckBits = conn.packetsIGot
	conn.mutex.RUnlock()
	packet.Type = t
	packet.Data = data
	buf := h.bufferPool.Get().([]byte)
	defer func() {
		buf = buf[:cap(buf)]
		h.bufferPool.Put(buf)
	}()
	fmt.Println("SENT packet:", packet)
	n, err := packet.Encode(buf)
	if err != nil {
		return err
	}
	// client server write separation
	if conn.Addr == nil {
		_, err = h.conn.Write(buf[:n])
	} else {
		_, err = h.conn.WriteToUDP(buf[:n], conn.Addr)
	}
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

// Give existed or create new connection
func (h *host) getConnection(key string, addr *net.UDPAddr) *connection {
	h.mu.RLock()
	conn, exist := h.connections[key]
	h.mu.RUnlock()
	if exist {
		return conn
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	conn, exist = h.connections[key]
	if exist {
		return conn
	}

	conn = &connection{
		state: state{
			id: 1,
		},
		packetsSend: make(map[uint32]SentPacket),
		Addr:        addr,
	}
	h.connections[key] = conn

	return conn
}

func (h *host) processRawPacket(conn *connection, data []byte, sendPacketHere chan<- Packet) error {
	packet := h.packetPool.Get().(*Packet)
	defer func() {
		*packet = Packet{}
		h.packetPool.Put(packet)
	}()
	packet.Addr = conn.Addr
	err := packet.Decode(data)
	if err != nil {
		return fmt.Errorf("transport:receive:decode: %w", err)
	}
	conn.mutex.Lock()
	conn.lastHeared = time.Now()
	isDuplicate := h.processPacket(conn, packet)
	conn.mutex.Unlock()
	// if packet double just drop
	// duplicate
	// continue
	if !isDuplicate {
		sendPacketHere <- *packet
	}
	return nil
}

// func (h *host) receive(conn *connection, sendPacketHere chan<- Packet) error {
// 	buf := h.bufferPool.Get().([]byte)
// 	defer func() {
// 		buf = buf[:cap(buf)]
// 		h.bufferPool.Put(buf)
// 	}()
// 	n, userAddr, err := h.conn.ReadFromUDP(buf)
// 	if err != nil {
// 		return fmt.Errorf("transport:receive:%w", err)
// 	}
//
// 	packet := h.packetPool.Get().(*Packet)
// 	defer func() {
// 		*packet = Packet{}
// 		h.packetPool.Put(packet)
// 	}()
// 	packet.Addr = userAddr
// 	err = packet.Decode(buf[:n])
// 	if err != nil {
// 		return fmt.Errorf("transport:receive:decode: %w", err)
// 	}
// 	conn.mutex.Lock()
// 	conn.lastHeared = time.Now()
// 	isDuplicate := h.processPacket(conn, packet)
// 	conn.mutex.Unlock()
// 	// if packet double just drop
// 	// duplicate
// 	// continue
// 	if !isDuplicate {
// 		sendPacketHere <- *packet
// 	}
// 	return nil
// }

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

func newHost(conn *net.UDPConn) *host {
	return &host{
		conn:        conn,
		connections: make(map[string]*connection),
		packetPool: sync.Pool{
			New: func() any {
				return &Packet{}
			},
		},
		bufferPool: sync.Pool{
			New: func() any {
				return make([]byte, bufferSize)
			},
		},
	}
}

func NewClient(h, port string) (*clientHost, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(h, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	clntHost := &clientHost{
		newHost(conn),
	}
	return clntHost, nil
}

func NewServer(h, port string) (*serverHost, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(h, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	srvHost := &serverHost{
		newHost(conn),
	}
	return srvHost, nil
}
