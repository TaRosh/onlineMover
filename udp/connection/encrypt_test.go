package connection

import (
	"testing"

	"github.com/TaRosh/online_mover/udp/packet"
	"github.com/stretchr/testify/require"
)

func TestEncryption(t *testing.T) {
	conn := New(nil)

	// TEST: encrypt without encryptor
	p := newPacket(1)
	p.Data = []byte("test")
	buf := make([]byte, 1024)
	n, err := p.Encode(buf)
	require.NoError(t, err)
	_, err = conn.EncryptPacket(buf[:n], &p)
	require.Error(t, err)

	// TEST: buffer small to encrypt
	_, err = conn.EncryptPacket(buf[:n-1], &p)
	require.Error(t, err)

	// TEST: encryption work
	conn.CreateEncryptor(testKey)
	p.Data = []byte("test")
	buf = make([]byte, 1024)
	n, _ = p.Encode(buf)
	var originalBytes []byte
	originalBytes = append(originalBytes, buf[:n]...)
	n, err = conn.EncryptPacket(buf, &p)
	require.NoError(t, err)
	require.Equal(t, packet.HeaderSize+len(p.Data)+conn.gcm.Overhead(), n)
	require.NotEqual(t, originalBytes, buf[:n])
	require.NotEqual(t, p.Data, buf[packet.HeaderSize:])
}
