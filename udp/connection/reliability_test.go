package connection

import (
	"testing"

	"github.com/TaRosh/online_mover/udp/packet"
	"github.com/stretchr/testify/require"
)

func resetConn(conn *Conn) {
	conn.seq = 0
	conn.lastSeqReceived = 0
	conn.packetsIGot = 0
}

func TestReliabilityConnection(t *testing.T) {
	// TEST: missing packet in ackbits
	conn := New(nil)
	require.NotNil(t, conn)

	var receivedPacket []packet.Packet
	receivedPacket = append(receivedPacket, newPacket(96))
	receivedPacket = append(receivedPacket, newPacket(98))
	receivedPacket = append(receivedPacket, newPacket(99))
	receivedPacket = append(receivedPacket, newPacket(100))
	for _, p := range receivedPacket {
		conn.ProcessHeaders(p.Header)
	}
	require.Equal(t, uint64(100), conn.lastSeqReceived)
	require.Equal(t, uint32(0b10111), conn.packetsIGot)

	// Test: ackbits when old packet arrive
	resetConn(conn)
	conn.lastSeqReceived = 100
	conn.packetsIGot = 0b1011
	p := newPacket(98)
	conn.ProcessHeaders(p.Header)
	require.Equal(t, uint32(0b1111), conn.packetsIGot)
	require.Equal(t, uint64(100), conn.lastSeqReceived)

	// Test: 32 history boundary
	resetConn(conn)
	receivedPacket = receivedPacket[:0]
	for i := uint64(100); i > 60; i-- {
		receivedPacket = append(receivedPacket, newPacket(i))
	}
	for _, p := range receivedPacket {
		conn.ProcessHeaders(p.Header)
	}
	require.Equal(t, uint64(100), conn.lastSeqReceived)
	require.NotEqual(t, uint32(0), conn.packetsIGot)

	// Test: detect dublicates
	resetConn(conn)
	receivedPacket = receivedPacket[:0]
	receivedPacket = append(receivedPacket, newPacket(100))
	receivedPacket = append(receivedPacket, newPacket(100))
	isDuplicate := false
	isDuplicate = conn.ProcessHeaders(receivedPacket[0].Header)
	require.Equal(t, false, isDuplicate)
	isDuplicate = conn.ProcessHeaders(receivedPacket[1].Header)
	require.Equal(t, true, isDuplicate)

	// Test: sequence wraparound
	// as we use uint32: 4294967295 → 0
	// we need it becouse seq: 2 - 4294967295 =
	resetConn(conn)
	conn.lastSeqReceived = ^uint64(0)
	p = newPacket(0)
	conn.ProcessHeaders(p.Header)
	require.Equal(t, uint64(0), conn.lastSeqReceived)

	// Test: Recive packet with seq bigger then our window(ackbits length)
	resetConn(conn)
	conn.lastSeqReceived = 0
	p = newPacket(32)
	conn.ProcessHeaders(p.Header)
	require.Equal(t, uint32(1), conn.packetsIGot)
}
