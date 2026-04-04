package connection

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptionKeys(t *testing.T) {
	conn := New(nil)
	// TEST: generated shared key is valid
	pub, priv, shared, err := conn.GenerateSharedKey(testKey)
	require.NoError(t, err)
	require.NotNil(t, pub)
	require.NotNil(t, priv)
	require.NotNil(t, shared)

	// TEST: add keys work
	conn.AddKeys(pub, priv)
	require.Equal(t, pub, conn.publicKey)
	require.Equal(t, priv, conn.privateKey)
}
