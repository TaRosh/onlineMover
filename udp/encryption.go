package udp

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
)

func (h *host) sendEncryptionRequest(conn *connection) error {
	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("host:sendEncryptionRequest: %w", err)
	}
	conn.addKeys(priv.PublicKey().Bytes(), priv)
	fmt.Println("Priv len", len(priv.PublicKey().Bytes()))
	err = h.send(conn, PacketKeyExchangeRequest, priv.PublicKey().Bytes())
	if err != nil {
		return fmt.Errorf("host:sendEncryptionRequest: %w", err)
	}

	conn.changeEncryptionState(EncryptionSendRequest)
	return nil
}

// function on receive side. Get remote key, create encryption and send
// local public key to remote side
func (h *host) handleEncryptionRequest(conn *connection, remoteKey []byte) error {
	myPriv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("host:proccessRawPacket: %w", err)
	}
	publicKey := myPriv.PublicKey()
	conn.addKeys(publicKey.Bytes(), myPriv)
	remotePub, err := ecdh.X25519().NewPublicKey(remoteKey)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionRequest: %w", err)
	}
	shared, err := myPriv.ECDH(remotePub)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionRequest: %w", err)
	}
	err = conn.createEncryption(shared)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionRequest: %w", err)
	}
	// send my public key
	err = h.send(conn, PacketKeyExchangeAnswer, conn.getPublicKey())
	if err != nil {
		return fmt.Errorf("host:handleEncryptionRequest: %w", err)
	}
	// now server should now about key and we can send encrypted data
	return nil
}

// Function for receving remote key for our encryption hello request.
func (h *host) handleEncryptionResponse(conn *connection, remoteKey []byte) error {
	remotePub, err := ecdh.X25519().NewPublicKey(remoteKey)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionResponse: %w", err)
	}
	priv := conn.getPrivateKey()
	if err != nil {
		return fmt.Errorf("host:handleEncryptionResponse: %w", err)
	}
	shared, err := priv.ECDH(remotePub)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionResponse: %w", err)
	}
	err = conn.createEncryption(shared)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionResponse: %w", err)
	}
	return nil
}
