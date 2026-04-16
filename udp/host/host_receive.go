package host

import (
	"fmt"
	"log"

	inner "github.com/TaRosh/online_mover/udp/connection"
	p "github.com/TaRosh/online_mover/udp/packet"
)

// receive data
// 1. Decode public header
// 2. Decode rest of the packet
// 3. Run handler based on packet type
// 4. Update connection based on packet header ( seq, ack, ackBits )
func (h *host) receive(conn *inner.Conn, data []byte, packet *p.Packet) error {
	var err error
	// var n int

	// if public header corrupted can't process packet -> return
	_, err = packet.PublicHeader.Decode(data)
	if err != nil {
		return fmt.Errorf("host:handleIncomingPacket: %w", err)
	}

	// 1. Decode packet
	// fmt.Printf("RECV packet: %+v\n", packet)

	err = h.handleConn[conn.GetEncryptionState()](conn, packet, data)
	// if packet corrupted we can't get full headers -> return
	if err != nil {
		return fmt.Errorf("host:handleIncomingPacket: %w", err)
	}
	// here packet decrypted
	// and we can run middleware
	// in our case it is function that need run over packet
	// just now it change addr in conn to addr from packet
	// It valid becous if connection unsecure -> trust it
	// If connection secure and we decrypt packet without err -> trust it
	for _, mwareFn := range h.middleware {

		err := mwareFn(conn, packet)
		if err != nil {
			return err
		}
	}

	// 2. Handle if it part of encryption handshake
	fn, exist := h.handleEncrypt[packet.Header.Type]
	if exist {
		err = fn(conn, packet)
		// here we just for some reason can' make connection secure
		// caller must check err and resend request if client
		// and wait or send err if server
		// just log and return err after proccessing headers
		log.Println("host:handleIncomingPacket:encryption handshak:", err)
		// dont update our reliability if it part of handshake
		return err
	}

	// 3. Reliability
	// here packet decoded and decrypted
	// fmt.Printf("RECV packet: %+v\n", packet)
	if conn.ProcessHeaders(packet.Header) {
		return fmt.Errorf("host:receive: %w", ErrDuplicatePacket)
	}
	return nil
}
