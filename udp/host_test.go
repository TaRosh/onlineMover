package udp

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func resetConn(conn *connection) {
	conn.state.id = 0
	conn.state.lastIDReceived = 0
	conn.state.packetsIGot = 0
}

func newPacket(seq uint32) Packet {
	return Packet{
		Header: Header{
			Sequence: seq,
		},
	}
}

func TestProcessPacket(t *testing.T) {
	// Test: missing packet in ackbits
	h, err := NewServer("localhost", "3000")
	require.NoError(t, err)
	require.NotNil(t, h)
	var receivedPacket []Packet
	receivedPacket = append(receivedPacket, newPacket(96))
	receivedPacket = append(receivedPacket, newPacket(98))
	receivedPacket = append(receivedPacket, newPacket(99))
	receivedPacket = append(receivedPacket, newPacket(100))
	conn := connection{
		state: state{
			id:             0,
			lastIDReceived: 0,
			packetsIGot:    0,
		},
	}
	for _, p := range receivedPacket {
		h.processPacket(&conn, &p)
	}
	require.Equal(t, uint32(100), conn.state.lastIDReceived)
	require.Equal(t, uint32(0b10111), conn.state.packetsIGot)

	// Test: ackbits when old packet arrive
	resetConn(&conn)
	conn.state.lastIDReceived = 100
	conn.state.packetsIGot = 0b1011
	packet := newPacket(98)
	h.processPacket(&conn, &packet)
	require.Equal(t, uint32(0b1111), conn.state.packetsIGot)
	require.Equal(t, uint32(100), conn.state.lastIDReceived)

	// Test: 32 history boundary
	resetConn(&conn)
	receivedPacket = receivedPacket[:0]
	for i := uint32(100); i > 60; i-- {
		receivedPacket = append(receivedPacket, newPacket(i))
	}
	for _, p := range receivedPacket {
		h.processPacket(&conn, &p)
	}
	require.Equal(t, uint32(100), conn.lastIDReceived)
	require.NotEqual(t, uint32(0), conn.packetsIGot)

	// Test: detect dublicates
	resetConn(&conn)
	receivedPacket = receivedPacket[:0]
	receivedPacket = append(receivedPacket, newPacket(100))
	receivedPacket = append(receivedPacket, newPacket(100))
	isDuplicate := false
	for _, p := range receivedPacket {
		isDuplicate = h.processPacket(&conn, &p)
	}
	require.Equal(t, true, isDuplicate)

	// Test: sequence wraparound
	// as we use uint32: 4294967295 → 0
	// we need it becouse seq: 2 - 4294967295 =
	resetConn(&conn)
	conn.lastIDReceived = ^uint32(0)
	p := newPacket(0)
	h.processPacket(&conn, &p)
	require.Equal(t, uint32(0), conn.lastIDReceived)

	// Test: Recive packet with seq bigger then our window(ackbits length)
	resetConn(&conn)
	conn.lastIDReceived = 0
	p = newPacket(32)
	h.processPacket(&conn, &p)
	require.Equal(t, uint32(1), conn.packetsIGot)

	// Test:
}

func TestSentFunction(t *testing.T) {
	server, err := NewServer("localhost", "9001")
	require.NoError(t, err)
	require.NotNil(t, server)
	defer server.Close()

	client, err := NewClient("localhost", "9000")
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	packets := make(chan Packet, 10)

	go server.Receive(packets)

	err = client.Send(PacketConnect, nil)
	require.NoError(t, err)

	// Test:
}

func TestReceiveFunction(t *testing.T) {
	h, err := NewServer("localhost", "9000")
	require.NoError(t, err)
	require.NotNil(t, h)

	packets := make(chan Packet, 10)
	buf := make([]byte, 1024)

	client, err := net.DialUDP("udp", nil, h.conn.LocalAddr().(*net.UDPAddr))
	require.NoError(t, err)

	// Test: packet delivery detection
	p := newPacket(1)
	n, err := p.Encode(buf)
	require.NoError(t, err)

	client.Write(buf[:n])
	h.Receive(packets)
	packet := <-packets
	require.Equal(t, p.Header, packet.Header)
	require.Equal(t, p.Data, packet.Data)

	// Test: drop dublicates
	p = newPacket(2)
	n, err = p.Encode(buf)
	require.NoError(t, err)
	client.Write(buf[:n])
	h.Receive(packets)

	client.Write(buf[:n])
	h.Receive(packets)

	time.Sleep(100 * time.Millisecond)
	require.Equal(t, 1, len(packets))

	// Test: check seq ack ackbits from many client
}
