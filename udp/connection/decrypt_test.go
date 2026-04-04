package connection

import (
	"testing"

	"github.com/TaRosh/online_mover/udp/packet"
	"github.com/stretchr/testify/require"
)

func TestDecrypt(t *testing.T) {
	// TEST: decrypt packet
	conn := New(nil)
	buf := make([]byte, 1024)
	conn.CreateEncryptor(testKey)

	p := newPacket(1)
	p.Data = []byte("test")

	n, _ := p.Encode(buf)
	var originalBytes []byte
	originalBytes = append(originalBytes, buf[:n]...)

	n, err := conn.EncryptPacket(buf, &p)

	decryptedBytes, err := conn.Decrypt(p.PublicHeader, buf[:n])

	require.NoError(t, err)
	decryptedPacket := newPacket(0)
	n, _ = decryptedPacket.PrivateHeader.Decode(decryptedBytes)
	require.Equal(t, p.PrivateHeader, decryptedPacket.PrivateHeader)
	decryptedPacket.DecodeBody(decryptedBytes[n:])
	require.Equal(t, p.Data, decryptedPacket.Data)

	// TEST: try decrypt corrupted bytes

	buf[1] = buf[1] + 1
	decryptedBytes, err = conn.Decrypt(p.PublicHeader, buf)
	require.Error(t, err)

	buf[1] = buf[1] - 1
	buf[packet.HeaderSize+1] = buf[packet.HeaderSize+1] + 1
	decryptedBytes, err = conn.Decrypt(p.PublicHeader, buf)
	require.Error(t, err)
}
