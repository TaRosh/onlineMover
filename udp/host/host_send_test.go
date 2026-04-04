package host

import (
	"net"
	"testing"

	inner "github.com/TaRosh/online_mover/udp/connection"
	p "github.com/TaRosh/online_mover/udp/packet"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("localhost", "3000"))
	c, err := net.ListenUDP("udp", addr)
	hst := newHost(c)
	// require.NoError(t, err)
	require.NotNil(t, hst)
	defer hst.Close()

	// TEST: send raw packet on unsecure channel
	conn := inner.New(nil)
	data := []byte("test")
	buf := make([]byte, 1024)
	var n int
	err = hst.send(conn, p.Connect, data, func(b []byte, u *net.UDPAddr) error {
		n = len(b)
		buf = append(buf[:0], b...)
		return nil
	})
	require.NoError(t, err)

	packet := newPacket(1)
	packet.Decode(buf[:n])
	require.Equal(t, data, packet.Data)

	// TODO: not here becouse we get id when server
	// first time receive msg from new addr
	// TEST: packet connect id should change for host new conn id
	// require.NotEqual(t, uint32(0), packet.ConnectionID)

	// TEST: receive packet
	// hst.handleIncomingPacket(conn, buf[:len(data)])

	// TEST: response packet have new connection id from host
	// packet := newPacket(1)

	// TEST: receive packet with client id
	//
	// p := <-packets
	// require.NotEqual(t, uint32(0), p.PublicHeader.ConnectionID)
	// require.NotNil(t, hst.connections[p.ConnectionID])
	// require.NotNil(t, hst.addrToConn[p.Addr.String()])
	//
	// require.NoError(t, err)

	// Test:
}
