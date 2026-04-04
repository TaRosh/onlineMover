package packet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeaderEncodeDecode(t *testing.T) {
	// TEST: small buffer
	original := Header{
		PublicHeader: PublicHeader{
			ConnectionID: 1,
			Sequence:     1,
		},
		PrivateHeader: PrivateHeader{
			Ack:     1,
			AckBits: 1,
			Type:    1,
		},
	}
	buf := make([]byte, HeaderSize-1)
	_, err := original.Encode(buf)
	require.Error(t, err)

	// TEST: encode

	buf = make([]byte, HeaderSize)
	n, err := original.Encode(buf)
	require.NoError(t, err)
	require.Equal(t, n, HeaderSize)

	//  TEST: to small data to decode
	var pHub Header
	n, err = pHub.Decode(buf[:n-1])
	require.Error(t, err)
	require.Equal(t, n, 0)

	// TEST: decode
	n, err = pHub.Decode(buf)
	require.NoError(t, err)
	require.Equal(t, n, HeaderSize)
	require.Equal(t, original, pHub)
}
