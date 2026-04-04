package connection

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
)

func (conn *Conn) AddKeys(publicKey []byte, privateKey *ecdh.PrivateKey) {
	// TODO: probably need add check if keys already seted
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.privateKey = privateKey
	conn.publicKey = publicKey
}

func (conn *Conn) GenerateSharedKey(remoteKey []byte) (publicKey []byte, privateKey *ecdh.PrivateKey, sharedKey []byte, err error) {
	privateKey, err = ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Conn:GenerateSharedKey: %w", err)
	}
	publicKey = privateKey.PublicKey().Bytes()
	remotePub, err := ecdh.X25519().NewPublicKey(remoteKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Conn:GenerateSharedKey: %w", err)
	}
	sharedKey, err = privateKey.ECDH(remotePub)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Conn:GenerateSharedKey: %w", err)
	}
	return publicKey, privateKey, sharedKey, err
}

func (conn *Conn) GetPublicKey() []byte {
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	return conn.publicKey
}

func (conn *Conn) GetPrivateKey() *ecdh.PrivateKey {
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	return conn.privateKey
}
