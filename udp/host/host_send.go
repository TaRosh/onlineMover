package host

import (
	"fmt"
	"net"
	"time"

	inner "github.com/TaRosh/online_mover/udp/connection"
	"github.com/TaRosh/online_mover/udp/packet"
)

// send data.
// 1. Set packet with connection state, packet type and data
// 2. Encrypt packet (private header + data ) if connection secure
// 3. Send raw packet bytes via writefn function
func (h *host) send(conn *inner.Conn, t packet.Type, data []byte, writeFn func([]byte, *net.UDPAddr) error) error {
	pack := h.getPacket()
	defer h.putPacket(pack)

	//*** SET PACKET
	connState := conn.GetConnectionState()
	// pub
	pack.PublicHeader.Sequence = connState.Seq
	pack.PublicHeader.ConnectionID = connState.ID

	// priv
	pack.PrivateHeader.Ack = connState.Ack
	pack.PrivateHeader.AckBits = connState.AckBits
	pack.PrivateHeader.Type = t

	// data
	pack.Data = data
	// **** END OF SET PACKET

	buf := h.getBuf()
	defer h.putBuf(buf)

	var err error
	var n int

	// fmt.Printf("SENT packet: %+v\n", pack)
	// send raw or encrypted packet
	if conn.GetEncryptionState() != inner.Secure {

		n, err = pack.Encode(buf)
		if err != nil {
			return fmt.Errorf("host:send: %w", err)
		}
	} else {

		n, err = conn.EncryptPacket(buf, pack)
		if err != nil {
			return fmt.Errorf("host:send: %w", err)
		}
	}

	// client/server write function
	err = writeFn(buf[:n], conn.Addr)
	if err != nil {
		return err
	}
	conn.ConsiderPacket(connState.Seq, packet.SentPacket{
		SendedWhen: time.Now(),
		Delivered:  false,
	})
	return nil
}
