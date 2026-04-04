package connection

import (
	"crypto/cipher"
	"crypto/ecdh"
	"net"
	"sync"
	"time"

	"github.com/TaRosh/online_mover/udp/packet"
)

type Conn struct {
	id uint32

	// network
	Addr *net.UDPAddr

	// reliability
	seq             uint64
	lastSeqReceived uint64
	packetsIGot     uint32

	// crypto
	state      SecureState
	gcm        cipher.AEAD
	privateKey *ecdh.PrivateKey
	publicKey  []byte

	// packte live
	packetsSend map[uint64]packet.SentPacket
	smoothedRTT time.Duration

	// where we last time receive packet from that connection
	lastHeared  time.Time
	lastCleanup time.Time
	mutex       sync.RWMutex
}

func (c *Conn) GetLastTimeActive() time.Time {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.lastHeared
}

func (c *Conn) SetID(id uint32) {
	c.mutex.Lock()
	c.id = id
	c.mutex.Unlock()
}

func (c *Conn) GetID() uint32 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.id
}

func (c *Conn) GetAddr() *net.UDPAddr {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Addr
}

func (c *Conn) SetAddr(playerAddr *net.UDPAddr) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Addr = playerAddr
}

func New(senderAddr *net.UDPAddr) *Conn {
	conn := Conn{
		state:       Unsecure,
		seq:         1,
		packetsSend: make(map[uint64]packet.SentPacket),
		Addr:        senderAddr,
		lastHeared:  time.Now(),
	}

	return &conn
}
