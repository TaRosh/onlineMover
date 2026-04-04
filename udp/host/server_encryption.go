package host

import (
	"fmt"

	inner "github.com/TaRosh/online_mover/udp/connection"
	"github.com/TaRosh/online_mover/udp/packet"
)

// TODO: think if error
// drop it becouse can't create secur?
// we get answer from server but can't
// create encryption on our side
// -> resend but in function upper
func (sh *serverHost) handleEncryptionRequest(conn *inner.Conn, pack *packet.Packet) error {
	// ignore duplicate
	if conn.GetEncryptionState() == inner.Secure {
		return nil
	}

	// 1. Create keys and encryptor
	remoteKey := pack.Data

	pub, priv, shared, err := conn.GenerateSharedKey(remoteKey)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionRequest: %w", err)
	}

	conn.AddKeys(pub, priv)
	// 2. Send my public key
	err = sh.host.send(conn, packet.KeyExchangeAnswer, pub, sh.write)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionRequest: %w", err)
	}
	err = conn.CreateEncryptor(shared)
	if err != nil {
		return fmt.Errorf("host:handleEncryptionRequest: %w", err)
	}

	// now server should now about key and we can send encrypted data
	return nil
}
