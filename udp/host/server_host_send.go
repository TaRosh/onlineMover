package host

import (
	"errors"
	"fmt"
	"net"

	"github.com/TaRosh/online_mover/udp/packet"
)

// responsible to write data in actual connection
func (sh *serverHost) write(data []byte, addr *net.UDPAddr) error {
	n, err := sh.host.conn.WriteToUDP(data, addr)
	if n != len(data) {
		return errors.New("serverHost:write: can't write full data")
	}
	if err != nil {
		return fmt.Errorf("serverHost:write: %w", err)
	}
	return nil
}

// send data via connection abstraction
// pass out write function to write data to wire
func (sh *serverHost) Send(id uint32, t packet.Type, data []byte) error {
	conn, exist := sh.connections[id]
	if !exist {
		return errors.New(fmt.Sprintf("serverHost:Send: no connection with id: %d", id))
	}
	err := sh.send(conn, t, data, sh.write)
	return err
}
