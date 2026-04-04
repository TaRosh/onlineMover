package connection

import (
	"fmt"

	"github.com/TaRosh/online_mover/udp/packet"
)

func (conn *Conn) Decrypt(publicHeaders packet.PublicHeader, data []byte) ([]byte, error) {
	nonce, err := conn.makeNonce(publicHeaders.Sequence)
	if err != nil {
		return nil, fmt.Errorf("conn:Decrypt: %w", err)
	}
	return conn.gcm.Open(nil, nonce, data[packet.PublicHeaderSize:], data[:packet.PublicHeaderSize])
}
