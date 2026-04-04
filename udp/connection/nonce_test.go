package connection

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNonce(t *testing.T) {
	// // TEST: nonce making without encryptor (gcm)

	seq := uint64(math.MaxUint64)
	buf := make([]byte, 12)
	conn := New(nil)
	conn.id = 0
	nonce, err := conn.makeNonce(seq)
	require.Error(t, err)
	require.Nil(t, nonce)

	// TEST: nonce making
	conn.CreateEncryptor(testKey)

	nonce, err = conn.makeNonce(seq)
	require.NoError(t, err)
	require.NotNil(t, nonce)

	target := uint64(math.MaxUint64)
	binary.BigEndian.PutUint64(buf[4:], target)
	require.Equal(t, conn.gcm.NonceSize(), len(nonce))
	require.Equal(t, buf, nonce)
}
