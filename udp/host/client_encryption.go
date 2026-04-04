package host

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"

	"github.com/TaRosh/online_mover/udp/connection"
	inner "github.com/TaRosh/online_mover/udp/connection"
	"github.com/TaRosh/online_mover/udp/packet"
)

func (ch *clientHost) sendEncryptionRequest() error {
	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("host:sendEncryptionRequest: %w", err)
	}
	ch.connection.AddKeys(priv.PublicKey().Bytes(), priv)
	fmt.Println("Priv len", len(priv.PublicKey().Bytes()))
	err = ch.host.send(ch.connection, packet.KeyExchangeRequest, priv.PublicKey().Bytes(), ch.write)
	if err != nil {
		return fmt.Errorf("host:sendEncryptionRequest: %w", err)
	}
	ch.connection.SetEncryptionState(inner.Wait)

	return nil
}

// TODO: think if err drop it becouse can't create secur?
func (ch *clientHost) handleEncryptionResponse(conn *connection.Conn, pack *packet.Packet) error {
	// handle only if connection wait for it
	if conn.GetEncryptionState() != inner.Wait {
		// unexpected -> ignore
		return nil
	}

	// TODO: add set connection id from server
	conn.SetID(pack.ConnectionID)

	remoteKey := pack.Data
	remotePub, err := ecdh.X25519().NewPublicKey(remoteKey)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionResponse: %w", err)
	}
	priv := conn.GetPrivateKey()
	if err != nil {
		return fmt.Errorf("host:handleEncryptionResponse: %w", err)
	}
	shared, err := priv.ECDH(remotePub)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionResponse: %w", err)
	}
	err = conn.CreateEncryptor(shared)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionResponse: %w", err)
	}
	// TODO: here state is secure so resend all queue packets
	return nil
}
