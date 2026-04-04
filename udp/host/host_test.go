package host

import (
	"net"
	"testing"

	"github.com/TaRosh/online_mover/udp/packet"
	"github.com/stretchr/testify/require"
)

func newPacket(seq uint64) packet.Packet {
	return packet.Packet{
		Header: packet.Header{
			PublicHeader: packet.PublicHeader{
				Sequence: seq,
			},
		},
	}
}

func TestNew(t *testing.T) {
	// TEST: create new host
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("localhost", "3000"))
	require.NoError(t, err)
	c, err := net.ListenUDP("udp", addr)
	require.NoError(t, err)
	h := newHost(c)
	defer h.Close()

	require.NotNil(t, h)
	require.Equal(t, c, h.conn)
	require.NotNil(t, h.handleConn)
	require.NotNil(t, h.handleEncrypt)
	require.NotNil(t, h.middleware)
	require.NotNil(t, h.packetPool.New)
	require.NotNil(t, h.bufferPool.New)
}
