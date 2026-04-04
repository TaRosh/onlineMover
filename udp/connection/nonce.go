package connection

import (
	"encoding/binary"
	"fmt"
)

func (conn *Conn) makeNonce(seq uint64) ([]byte, error) {
	if conn.gcm == nil {
		return nil, fmt.Errorf("conn:makeNonce: %w", ErrNoEncryptor)
	}
	nonce := make([]byte, conn.gcm.NonceSize())
	binary.BigEndian.PutUint32(nonce[0:4], conn.id)
	binary.BigEndian.PutUint64(nonce[4:12], seq)
	return nonce, nil
}
