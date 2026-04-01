package udp

import "net"

type clientHost struct {
	*host
	connection *connection
}

func (ch *clientHost) Send(t packetType, data []byte) error {
	// TODO: for now drop packets that insecure
	if ch.connection.encryptionState == EncryptionUnsecure {
		err := ch.sendEncryptionRequest(ch.connection)
		return err
	}
	// TODO: queue pending packets
	err := ch.send(ch.connection, t, data)
	return err
}

func (ch *clientHost) Receive(sendPacketHere chan<- Packet) error {
	buf := ch.bufferPool.Get().([]byte)
	defer func() {
		buf = buf[:cap(buf)]
		ch.bufferPool.Put(buf)
	}()
	n, _, err := ch.conn.ReadFromUDP(buf)
	err = ch.processRawPacket(ch.connection, buf[:n], sendPacketHere)
	if err != nil {
		return err
	}
	return nil
}

func NewClient(h, port string) (*clientHost, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(h, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	clntHost := &clientHost{
		host:       newHost(conn),
		connection: newConnection(nil),
	}
	return clntHost, nil
}
