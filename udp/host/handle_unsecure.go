package host

import (
	inner "github.com/TaRosh/online_mover/udp/connection"
	p "github.com/TaRosh/online_mover/udp/packet"
)

// public headers decoded before in host receive function
func handleUnsecure(conn *inner.Conn, packet *p.Packet, data []byte) error {
	// decode rest of headers
	data = data[p.PublicHeaderSize:]
	n, err := packet.PrivateHeader.Decode(data)
	if err != nil {
		return err
	}

	// decode rest
	packet.DecodeBody(data[n:])
	return nil
}
