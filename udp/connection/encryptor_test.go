package connection

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptor(t *testing.T) {
	conn := New(nil)
	// TEST: add encryptor work
	err := conn.CreateEncryptor(testKey)
	require.NoError(t, err)
	require.NotNil(t, conn.gcm)
}
