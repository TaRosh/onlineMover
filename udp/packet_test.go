package udp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPacketEncodeDecode(t *testing.T) {
	// we skip addr for packet becouse it set inside server
	// Test: equality with data
	header := Header{
		PublicHeader: PublicHeader{
			ConnectionID: 0,
			SeqShort:     0,
		},
		PrivateHeader: PrivateHeader{
			Sequence: 10,
			Ack:      8,
			AckBits:  0b11111111,
			Type:     PacketConnect,
		},
	}
	original := Packet{
		Header: header,
		Data:   []byte{1, 2, 3, 4},
	}
	buf := make([]byte, 1024)
	n, err := original.Encode(buf)
	require.NoError(t, err)
	require.NotEqual(t, 0, n)
	decoded := Packet{}
	err = decoded.Decode(buf[:n])
	require.NoError(t, err)
	require.Equal(t, original, decoded)

	// Test: equality with no data field
	header = Header{
		PrivateHeader: PrivateHeader{
			Sequence: 10,
			Ack:      8,
			AckBits:  0b11111111,
			Type:     PacketConnect,
		},
	}
	original = Packet{
		Header: header,
	}
	buf = make([]byte, 1024)
	n, err = original.Encode(buf)
	require.NoError(t, err)
	require.NotEqual(t, 0, n)
	decoded = Packet{}
	err = decoded.Decode(buf[:n])
	require.NoError(t, err)
	require.Equal(t, original, decoded)
}
