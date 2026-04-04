package packet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublicHeaderEncodeDecode(t *testing.T) {
	// TEST: small buffer
	original := PublicHeader{
		ConnectionID: 11,
		Sequence:     21,
	}
	buf := make([]byte, PublicHeaderSize-1)
	_, err := original.Encode(buf)
	require.Error(t, err)

	// TEST: encode

	buf = make([]byte, PublicHeaderSize)
	n, err := original.Encode(buf)
	require.NoError(t, err)
	require.Equal(t, n, PublicHeaderSize)

	//  TEST: to small data to decode
	var pHub PublicHeader
	n, err = pHub.Decode(buf[:n-1])
	require.Error(t, err)
	require.Equal(t, n, 0)

	// TEST: decode
	n, err = pHub.Decode(buf)
	require.NoError(t, err)
	require.Equal(t, n, PublicHeaderSize)
	require.Equal(t, original, pHub)
}
