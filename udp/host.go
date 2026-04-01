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

type host struct {
	conn *net.UDPConn

	bufferPool sync.Pool
	packetPool sync.Pool

	// sentPacket    *Packet
	// receivePacket *Packet
}

func (h *host) GetAddr() *net.UDPAddr {
	return h.conn.LocalAddr().(*net.UDPAddr)
}

// send encrypted data
func (h *host) send(conn *connection, t packetType, data []byte) error {
	packet := h.packetPool.Get().(*Packet)
	defer func() {
		*packet = Packet{}
		h.packetPool.Put(packet)
	}()

	//*** SET PACKET
	conn.mutex.RLock()
	packetState := conn.packetState
	packet.PrivateHeader.Sequence = packetState.id
	packet.PrivateHeader.Ack = packetState.lastIDReceived
	packet.PrivateHeader.AckBits = packetState.packetsIGot

	packet.PublicHeader.ConnectionID = conn.id
	conn.mutex.RUnlock()
	packet.PrivateHeader.Type = t
	packet.PublicHeader.SeqShort = uint16(packet.PrivateHeader.Sequence)

	packet.Data = data
	// **** END OF SET PACKET

	buf := h.bufferPool.Get().([]byte)
	defer func() {
		buf = buf[:cap(buf)]
		h.bufferPool.Put(buf)
	}()
	fmt.Printf("SENT packet: %+v\n", packet)
	var err error
	var n int
	if conn.getEncryptionState() != EncryptionSecure {
		// send raw
		// n = h.encodePacket
		n, err = packet.Encode(buf)
		if err != nil {
			return fmt.Errorf("host:send: %w", err)
		}
	} else {
		// encrypt packet

		n, err = conn.encryptPacket(buf, packet)
		if err != nil {
			return fmt.Errorf("host:send: %w", err)
		}
	}
	fmt.Println(n)

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

func (h *host) processRawPacket(conn *connection, data []byte, sendPacketHere chan<- Packet) error {
	packet := h.packetPool.Get().(*Packet)
	defer func() {
		*packet = Packet{}
		h.packetPool.Put(packet)
	}()
	packet.Addr = conn.Addr
	// TODO: add key exchange start
	// check key request or answer
	_, err := packet.Header.Decode(data[:headerSize])
	if err != nil {
		return fmt.Errorf("host:proccessRawPacket: %w", err)
	}
	if packet.PublicHeader.ConnectionID == 0 {
		packet.PublicHeader.ConnectionID = conn.id
	} else {
		conn.id = packet.PublicHeader.ConnectionID
	}

	fmt.Printf("Connection: %+v\n", conn)
	switch conn.encryptionState {
	case EncryptionUnsecure:
		packet.Data = data[headerSize:]

	case EncryptionSecure:
		// decrypt data first
		nonce := conn.makeNonce(packet.Header.Sequence)
		plainText, err := conn.decrypt(data[:headerSize], data[headerSize:], nonce)
		if err != nil {
			return fmt.Errorf("host:proccessRawPacket: %w", err)
		}
		packet.Data = plainText

	}
	// we server end get public key from client
	switch packet.Header.Type {
	case PacketKeyExchangeRequest:
		// packet.Data is a publicKey of the opposite side
		// generate private and public key for this connection on my side
		// use packet.Data ( public key ) to make shared key for connection
		// and from shared i create gcm
		err := h.handleEncryptionRequest(conn, packet.Data)
		if err != nil {
			return fmt.Errorf("host:proccessRawPacket: %w", err)
		}
		// TODO: think about not return here becouse after we process
		// seq, ack etc bits from packet header
	case PacketKeyExchangeAnswer:
		// we client and get public key from server
		err := h.handleEncryptionResponse(conn, packet.Data)
		if err != nil {
			return fmt.Errorf("host:proccessRawPacket: %w", err)
		}
	}
	// if packet not part encryption exchange
	// then it is raw packet ( if connection not secure )
	// or encrypted packet ( if connection secure )

	// Here packet should be decoded and encrypted

	fmt.Printf("RECV packet: %+v\n", packet)

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
		conn: conn,
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
