package host

import (
	"fmt"

	inner "github.com/TaRosh/online_mover/udp/connection"
	p "github.com/TaRosh/online_mover/udp/packet"
)

func handleSecure(conn *inner.Conn, packet *p.Packet, data []byte) error {
	// 2. Decrypt
	plaintext, err := conn.Decrypt(packet.PublicHeader, data)
	if err != nil {
		// drop it becouse can track it
		return fmt.Errorf("host:handleIncomingPacket: %w", err)
	}
	// decode the rest of the packet
	// private headers and body
	n, err := packet.PrivateHeader.Decode(plaintext)
	if err != nil {
		// can't get full headers -> return
		return err
	}
	packet.DecodeBody(plaintext[n:])
	return nil
}
