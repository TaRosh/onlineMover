package connection

import (
	"errors"
	"fmt"

	"github.com/TaRosh/online_mover/udp/packet"
)

var ErrNoEncryptor = errors.New("use encryption without encryptor (gcm)")

func (conn *Conn) EncryptPacket(buf []byte, pack *packet.Packet) (int, error) {
	if conn.GetEncryptionState() == Unsecure {
		return 0, fmt.Errorf("connection:encryptPacket: %w", ErrNoEncryptor)
	}
	if len(buf) < packet.HeaderSize+len(pack.Data)+conn.gcm.Overhead() {
		return 0, errors.New("connection:encryptPacket: buffer is to small for this data ( header + data + encrypt overhead)")
	}
	n, err := pack.PublicHeader.Encode(buf)
	if err != nil {
		return 0, fmt.Errorf("conn:encryptPacket: %w", err)
	}
	// buf[publicheader]
	pubHeadersEnd := n
	plaintextLen, err := pack.PrivateHeader.Encode(buf[pubHeadersEnd:])
	if err != nil {
		return 0, fmt.Errorf("conn:encryptPacket: %w", err)
	}

	// buf[public header|private header]
	copy(buf[pubHeadersEnd+plaintextLen:], pack.Data)
	// buf [public header | private header | data]
	//      pubHeadersEnd ^         plaintextLen ^
	plaintextLen += len(pack.Data)

	nonce, err := conn.makeNonce(pack.PublicHeader.Sequence)
	if err != nil {
		return 0, fmt.Errorf("conn:EncryptPacket: %w", err)
	}

	// public header + encrypt(priv heder + data)
	encryptedBuf := conn.encrypt(buf[:pubHeadersEnd], buf[:pubHeadersEnd], buf[pubHeadersEnd:pubHeadersEnd+plaintextLen], nonce)
	return len(encryptedBuf), nil
}

func (conn *Conn) encrypt(buf, publicHeaders, plainText, nonce []byte) []byte {
	return conn.gcm.Seal(buf, nonce, plainText, publicHeaders)
}
