package udp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"time"
)

type packetState struct {
	id             uint32
	lastIDReceived uint32
	packetsIGot    uint32
}

type secureState int

const (
	EncryptionUnsecure secureState = iota
	EncryptionSendRequest
	EncryptionReceiveResponse
	EncryptionSecure
)

type connection struct {
	id uint32
	packetState
	encryptionState secureState
	Addr            *net.UDPAddr
	packetsSend     map[uint32]SentPacket
	smoothedRTT     time.Duration

	// where we last time receive packet from that connection
	lastHeared  time.Time
	lastCleanup time.Time
	mutex       sync.RWMutex

	// Encryption level
	gcm        cipher.AEAD
	privateKey *ecdh.PrivateKey
	publicKey  []byte
}

func (c *connection) createEncryption(sharedKey []byte) error {
	block, err := aes.NewCipher(sharedKey)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.gcm = gcm
	c.encryptionState = EncryptionSecure
	return nil
}

func (c *connection) encryptPacket(buf []byte, packet *Packet) (int, error) {
	if c.encryptionState != EncryptionSecure {
		return 0, errors.New("connection:encryptPacket: use encryption without gcm be setted")
	}
	if len(buf) < headerSize+len(packet.Data)+c.gcm.Overhead() {
		return 0, errors.New("connection:encryptPacket: buffer is to small for this data ( header + data + encrypt overhead)")
	}
	n := packet.PublicHeader.Encode(buf)
	// buf[publicheader]
	pubHeadersEnd := n
	plaintextLen := packet.PrivateHeader.Encode(buf[pubHeadersEnd:])

	// buf[public header|private header]
	copy(buf[pubHeadersEnd+plaintextLen:], packet.Data)
	// buf [public header | private header | data]
	//      pubHeadersEnd ^         plaintextLen ^
	plaintextLen += len(packet.Data)

	nonce := c.makeNonce(packet.PrivateHeader.Sequence)

	// public header + encrypt(priv heder + data)
	encryptedBuf := c.encrypt(buf[:pubHeadersEnd], buf[:pubHeadersEnd], buf[pubHeadersEnd:pubHeadersEnd+plaintextLen], nonce)
	return len(encryptedBuf), nil
}

func (c *connection) encrypt(buf, publicHeaders, plainText, nonce []byte) []byte {
	return c.gcm.Seal(buf, nonce, plainText, publicHeaders)
}

func (c *connection) decrypt(publicHeaders, cipherText, nonce []byte) ([]byte, error) {
	return c.gcm.Open(nil, nonce, cipherText, publicHeaders)
}

func (c *connection) makeNonce(seq uint32) []byte {
	nonce := make([]byte, c.gcm.NonceSize())
	binary.BigEndian.PutUint32(nonce[0:4], c.id)
	binary.BigEndian.PutUint32(nonce[4:8], seq)
	offset := c.gcm.NonceSize() - 4
	binary.BigEndian.PutUint32(nonce[offset:], seq)
	return nonce
}

func (c *connection) addKeys(publicKey []byte, privateKey *ecdh.PrivateKey) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.privateKey = privateKey
	c.publicKey = publicKey
}

func (c *connection) getPublicKey() []byte {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.publicKey
}

func (c *connection) getPrivateKey() *ecdh.PrivateKey {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.privateKey
}

func (c *connection) changeEncryptionState(state secureState) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.encryptionState = state
}

func (c *connection) getEncryptionState() secureState {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.encryptionState
}

func newConnection(senderAddr *net.UDPAddr) *connection {
	conn := connection{
		packetState: packetState{
			id: 1,
		},
		encryptionState: EncryptionUnsecure,
		packetsSend:     make(map[uint32]SentPacket),
		Addr:            senderAddr,
	}

	return &conn
}
