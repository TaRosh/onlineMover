package connection

import (
	"time"

	"github.com/TaRosh/online_mover/udp/packet"
)

type ConnectionState struct {
	Seq     uint64
	Ack     uint64
	AckBits uint32
	ID      uint32
}

func (c *Conn) GetConnectionState() ConnectionState {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	cState := ConnectionState{
		Seq:     c.seq,
		Ack:     c.lastSeqReceived,
		AckBits: c.packetsIGot,
		ID:      c.id,
	}
	return cState
}

// check is packet id is old or new
func isNewer(a, b uint64) bool {
	return int64(a-b) > 0
}

// return true if it's duplicate packet
func (conn *Conn) ProcessHeaders(header packet.Header) (isDublicatePacket bool) {
	// last time heared here
	conn.lastHeared = time.Now()
	// 1. process ack

	conn.processAck(header)
	// 2. process
	// 2. process other
	// if packet id is new ( < then last id i got)
	// check is inside our check table ( uint32 ) or we lost
	// 32 packet and it is newer
	// do new check table or shift for it's diff ( currentId - lastId )
	if isNewer(header.Sequence, conn.lastSeqReceived) {
		shift := header.Sequence - conn.lastSeqReceived
		if shift > 32 {
			conn.packetsIGot = 0
		} else {
			conn.packetsIGot <<= uint32(shift)
		}
		conn.packetsIGot |= 1
		conn.lastSeqReceived = header.Sequence
	} else {
		diff := conn.lastSeqReceived - header.Sequence
		// TODO: packet not our window
		if diff > 32 && diff < 0 {
			return false
		}
		bitIndex := diff
		if conn.packetsIGot>>uint32(bitIndex)&1 != 0 {
			return true
		}
		conn.packetsIGot |= (1 << bitIndex)
	}
	return false
}

// mark packet as delivered and calculate rtt
func (conn *Conn) processAck(header packet.Header) {
	// check ack & ack bits for get rtt and calculate average rtt
	if sentPacket, ok := conn.getSentPacket(header.Ack); ok {
		if !sentPacket.Delivered {
			rawRtt := time.Since(sentPacket.SendedWhen)
			conn.calcRtt(rawRtt)
			conn.deleteSentPacket(header.Ack)
		}
	}
	for i := range uint64(32) {
		if (header.AckBits>>i)&1 == 1 {
			if sentPacket, ok := conn.getSentPacket(header.Ack - i); ok {
				rawRtt := time.Since(sentPacket.SendedWhen)
				conn.calcRtt(rawRtt)
				conn.deleteSentPacket(header.Ack - i)

			}
		}
	}
	// clean up the map for sended packets
	// delete data about packet if  not arrived after 3 second
	// launch every second
	if time.Since(conn.lastCleanup) > time.Second {
		for id, sentPacket := range conn.packetsSend {
			if time.Since(sentPacket.SendedWhen) > time.Second*3 {
				conn.deleteSentPacket(id)
			}
		}
		conn.lastCleanup = time.Now()
	}
}

func (conn *Conn) getSentPacket(id uint64) (packet.SentPacket, bool) {
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	p, exist := conn.packetsSend[id]

	return p, exist
}

// Add packet to rtt map and increment seq counter
func (conn *Conn) ConsiderPacket(seq uint64, sp packet.SentPacket) {
	conn.mutex.Lock()
	conn.packetsSend[seq] = sp
	conn.seq += 1
	conn.mutex.Unlock()
}

func (conn *Conn) calcRtt(sample time.Duration) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if conn.smoothedRTT == 0 {
		conn.smoothedRTT = sample
	} else {
		conn.smoothedRTT = time.Duration(float64(conn.smoothedRTT)*0.9 + float64(sample)*0.1)
	}
}

func (conn *Conn) deleteSentPacket(id uint64) {
	conn.mutex.Lock()
	delete(conn.packetsSend, id)
	conn.mutex.Unlock()
}
