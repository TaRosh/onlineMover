package packet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrivateHeaderEncodeDecode(t *testing.T) {
	// TEST: small buffer
	original := PrivateHeader{
		Ack:     10,
		AckBits: 1111,
		Type:    1,
	}
	buf := make([]byte, PublicHeaderSize-1)
	_, err := original.Encode(buf)
	require.Error(t, err)

	// TEST: encode

	buf = make([]byte, PrivateHeaderSize)
	n, err := original.Encode(buf)
	require.NoError(t, err)
	require.Equal(t, n, PrivateHeaderSize)

	//  TEST: to small data to decode
	var pHub PrivateHeader
	n, err = pHub.Decode(buf[:n-1])
	require.Error(t, err)
	require.Equal(t, n, 0)

	// TEST: decode
	n, err = pHub.Decode(buf)
	require.NoError(t, err)
	require.Equal(t, n, PrivateHeaderSize)
	require.Equal(t, original, pHub)
}
