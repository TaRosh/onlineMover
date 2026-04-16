package host

import (
	"errors"
	"fmt"
	"net"
	"time"

	inner "github.com/TaRosh/online_mover/udp/connection"
	"github.com/TaRosh/online_mover/udp/packet"
)

// responsible to write data in actual connection
func (ch *clientHost) write(data []byte, addr *net.UDPAddr) error {
	n, err := ch.host.conn.Write(data)
	if n != len(data) {
		return errors.New("clientHost:write: can't write full data")
	}
	if err != nil {
		return fmt.Errorf("clientHost:write: %w", err)
	}
	return nil
}

// send data via connection abstraction
// pass out write function to write data to wire
func (ch *clientHost) Send(t packet.Type, data []byte) error {
	var err error
	// initate hanshake
	// TODO: add lost packet to queue
	switch ch.connection.GetEncryptionState() {
	case inner.Unsecure:
		// TODO: add packet to queue
		err = ch.sendEncryptionRequest()
		return err
	case inner.Wait:
		// If still not secure -> resend encrypt request
		if time.Since(ch.lastEncryptionRequestSend) > 500*time.Millisecond {
			err = ch.sendEncryptionRequest()
			// TODO: add packet to queue
			return err
		}

	}
	// TODO: don't forgot send all packet from queue
	// when connection change state to secure
	err = ch.send(ch.connection, t, data, ch.write)
	return err
}
