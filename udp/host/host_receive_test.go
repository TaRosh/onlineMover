package host

import (
	"testing"

	"github.com/TaRosh/online_mover/udp/connection"
	"github.com/TaRosh/online_mover/udp/packet"
	"github.com/stretchr/testify/require"
)

func TestHostReceive(t *testing.T) {
	h := newHost(nil)
	conn := connection.New(nil)
	original := newPacket(1)
	original.Type = packet.KeyExchangeAnswer
	original.Data = []byte("test key")
	buf := make([]byte, 1024)

	n, _ := original.Encode(buf)
	target := newPacket(0)
	h.receive(conn, buf[:n], &target)
	require.NotEqual(t, original, target)
}
