package host

import "github.com/TaRosh/online_mover/udp/packet"

func (ch *clientHost) Receive(sendPacketHere chan<- packet.Packet) error {
	buf := ch.getBuf()
	defer ch.putBuf(buf)

	n, _, err := ch.conn.ReadFromUDP(buf)

	packet := ch.getPacket()
	defer ch.putPacket(packet)

	err = ch.host.receive(ch.connection, buf[:n], packet)
	if err != nil {
		return err
	}

	sendPacketHere <- *packet

	return nil
}
