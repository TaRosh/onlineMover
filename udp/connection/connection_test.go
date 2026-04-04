package connection

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

var testKey = []byte("12345678901234567890123456789012")

func TestNewConnection(t *testing.T) {
	// TEST: new connection give right seq number
	// & set addr
	addr, _ := net.ResolveUDPAddr("udp", "localhost:3000")
	conn := New(addr)
	require.Equal(t, uint64(1), conn.seq)
	require.Equal(t, addr, conn.Addr)
}
