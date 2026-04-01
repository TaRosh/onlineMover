package udp

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func resetConn(conn *connection) {
	conn.packetState.id = 1
	conn.packetState.lastIDReceived = 0
	conn.packetState.packetsIGot = 0
}

func newPacket(seq uint32) Packet {
	return Packet{
		Header: Header{
			PrivateHeader: PrivateHeader{
				Sequence: seq,
			},
		},
	}
}

var testKey = []byte("12345678901234567890123456789012")

func TestProcessPacket(t *testing.T) {
	// Test: missing packet in ackbits
	h, err := NewServer("localhost", "3000")
	require.NoError(t, err)
	require.NotNil(t, h)
	defer h.Close()

	var receivedPacket []Packet
	receivedPacket = append(receivedPacket, newPacket(96))
	receivedPacket = append(receivedPacket, newPacket(98))
	receivedPacket = append(receivedPacket, newPacket(99))
	receivedPacket = append(receivedPacket, newPacket(100))
	conn := connection{
		packetState: packetState{
			id:             0,
			lastIDReceived: 0,
			packetsIGot:    0,
		},
	}
	for _, p := range receivedPacket {
		h.processPacket(&conn, &p)
	}
	require.Equal(t, uint32(100), conn.packetState.lastIDReceived)
	require.Equal(t, uint32(0b10111), conn.packetState.packetsIGot)

	// Test: ackbits when old packet arrive
	resetConn(&conn)
	conn.packetState.lastIDReceived = 100
	conn.packetState.packetsIGot = 0b1011
	packet := newPacket(98)
	h.processPacket(&conn, &packet)
	require.Equal(t, uint32(0b1111), conn.packetState.packetsIGot)
	require.Equal(t, uint32(100), conn.packetState.lastIDReceived)

	// Test: 32 history boundary
	resetConn(&conn)
	receivedPacket = receivedPacket[:0]
	for i := uint32(100); i > 60; i-- {
		receivedPacket = append(receivedPacket, newPacket(i))
	}
	for _, p := range receivedPacket {
		h.processPacket(&conn, &p)
	}
	require.Equal(t, uint32(100), conn.lastIDReceived)
	require.NotEqual(t, uint32(0), conn.packetsIGot)

	// Test: detect dublicates
	resetConn(&conn)
	receivedPacket = receivedPacket[:0]
	receivedPacket = append(receivedPacket, newPacket(100))
	receivedPacket = append(receivedPacket, newPacket(100))
	isDuplicate := false
	for _, p := range receivedPacket {
		isDuplicate = h.processPacket(&conn, &p)
	}
	require.Equal(t, true, isDuplicate)

	// Test: sequence wraparound
	// as we use uint32: 4294967295 → 0
	// we need it becouse seq: 2 - 4294967295 =
	resetConn(&conn)
	conn.lastIDReceived = ^uint32(0)
	p := newPacket(0)
	h.processPacket(&conn, &p)
	require.Equal(t, uint32(0), conn.lastIDReceived)

	// Test: Recive packet with seq bigger then our window(ackbits length)
	resetConn(&conn)
	conn.lastIDReceived = 0
	p = newPacket(32)
	h.processPacket(&conn, &p)
	require.Equal(t, uint32(1), conn.packetsIGot)

	// Test:
}

func TestHostSendFunction(t *testing.T) {
	server, err := NewServer("localhost", "9000")
	require.NoError(t, err)
	require.NotNil(t, server)
	defer server.Close()

	client, err := NewClient("localhost", "9000")
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	packets := make(chan Packet, 10)

	// TEST: send raw packet
	err = client.send(client.connection, PacketConnect, nil)
	require.NoError(t, err)

	err = server.Receive(packets)
	require.NoError(t, err)

	// TEST: receive packet with client id

	p := <-packets
	require.NotEqual(t, uint32(0), p.PublicHeader.ConnectionID)
	require.NotNil(t, server.connections[p.ConnectionID])
	require.NotNil(t, server.addrToConn[p.Addr.String()])

	require.NoError(t, err)

	// Test:
}

func TestEncryptionLayer(t *testing.T) {
	// TEST: new connection is in unsecure state
	conn := newConnection(nil)
	require.Equal(t, EncryptionUnsecure, conn.encryptionState)
	require.Nil(t, conn.privateKey)
	require.Nil(t, conn.publicKey)

	// TEST: encrypt packet without gcm set

	packet := newPacket(1)
	packet.Data = []byte("test")
	buf := make([]byte, 1024)
	n, err := packet.Encode(buf)
	require.NoError(t, err)

	n, err = conn.encryptPacket(buf[:n], &packet)
	require.Error(t, err)

	// TEST: creat encryption gcm with shared key
	err = conn.createEncryption(testKey)
	require.NoError(t, err)
	require.NotNil(t, conn.gcm)

	// TEST: encrypt packet
	resetConn(conn)

	packet = newPacket(1)
	packet.Data = []byte("test")
	buf = make([]byte, 1024)
	n, err = packet.Encode(buf)
	var packetBytes []byte
	packetBytes = append(packetBytes, buf[:n]...)
	require.NoError(t, err)

	n, err = conn.encryptPacket(buf, &packet)
	require.NoError(t, err)
	require.Equal(t, headerSize+len(packet.Data)+conn.gcm.Overhead(), n)
	require.NotEqual(t, packetBytes, buf[:n])
	require.NotEqual(t, packet.Data, buf[headerSize:])

	// TEST: decrypt packet
	nonce := conn.makeNonce(packet.PrivateHeader.Sequence)
	m := packet.PublicHeader.Encode(buf)

	decryptedBuf, err := conn.decrypt(buf[:m], buf[m:n], nonce)
	require.NoError(t, err)
	require.Equal(t, packetBytes, append(buf[:m], decryptedBuf...))

	// TEST: test encryption handshake
	c, err := NewClient("localhost", "3000")
	require.NoError(t, err)
	require.NotNil(t, c)
	defer c.Close()

	s, err := NewServer("localhost", "3000")
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	err = c.sendEncryptionRequest(c.connection)
	require.NoError(t, err)
	require.Equal(t, EncryptionSendRequest, c.connection.encryptionState)
	require.NotNil(t, c.connection.privateKey)
	require.NotNil(t, c.connection.publicKey)
	packets := make(chan Packet, 10)

	serverConn := newConnection(c.conn.LocalAddr().(*net.UDPAddr))

	err = s.handleEncryptionRequest(serverConn, c.connection.publicKey)
	require.NoError(t, err)
	require.Equal(t, EncryptionSecure, serverConn.encryptionState)
	require.NotNil(t, serverConn.privateKey)
	require.NotNil(t, serverConn.publicKey)
	require.NotNil(t, serverConn.gcm)

	err = c.handleEncryptionResponse(c.connection, serverConn.publicKey)
	require.NoError(t, err)
	require.Equal(t, EncryptionSecure, c.connection.encryptionState)
	require.NotNil(t, c.connection.gcm)

	// TEST: send encrypted data
	resetConn(c.connection)
	resetConn(serverConn)

	data := []byte("test")

	err = c.Send(packet.Type, data)

	require.NoError(t, err)

	fmt.Printf("C: %+v\nS: %+v\n", c.connection, serverConn)
	s.addrToConn[c.conn.LocalAddr().String()] = serverConn
	err = s.Receive(packets)
	// err = s.processRawPacket(serverConn, buf[:n], packets)
	require.NoError(t, err)
	p := <-packets
	require.Equal(t, data, p.Data)

	// Test: make correct nonce
	nonce = c.connection.makeNonce(p.Header.Sequence)

	require.Equal(t, c.connection.gcm.NonceSize(), len(nonce))
	binary.BigEndian.PutUint32(buf, p.Sequence)
	require.Equal(t, nonce[c.connection.gcm.NonceSize()-4:], buf[:4])

	// Test: same data encrypted and decrypted

	cipherText := c.connection.encrypt(nil, nil, p.Data, nonce)
	require.NotEqual(t, cipherText, p.Data)

	plainText, err := c.connection.decrypt(nil, cipherText, nonce)
	require.NoError(t, err)
	require.Equal(t, plainText, p.Data)

	// Test: decrypt corrupted data
	cipherText[2] = 8
	plainText, err = c.connection.decrypt(nil, cipherText, nonce)
	require.Error(t, err)

	// Test: changed header data
	p = newPacket(2)
	p.Header.Ack = 11
	n, err = c.connection.encryptPacket(buf, &p)
	require.NoError(t, err)
	buf[3] = 3
	nonce = c.connection.makeNonce(p.Header.Sequence)
	_, err = c.connection.decrypt(buf[:headerSize], buf[headerSize:n], nonce)
	require.Error(t, err)
}

//
// func TestReceiveFunction(t *testing.T) {
// 	h, err := NewServer("localhost", "9000")
// 	require.NoError(t, err)
// 	require.NotNil(t, h)
// 	defer h.Close()
//
// 	packets := make(chan Packet, 10)
// 	buf := make([]byte, 1024)
//
// 	client, err := net.DialUDP("udp", nil, h.conn.LocalAddr().(*net.UDPAddr))
// 	require.NoError(t, err)
//
// 	// Test: packet delivery detection
// 	p := newPacket(1)
// 	p.Data = []byte("test")
// 	n, err := p.Encode(buf)
// 	require.NoError(t, err)
//
// 	client.Write(buf[:n])
// 	h.Receive(packets)
// 	packet := <-packets
// 	require.Equal(t, p.Header, packet.Header)
// 	require.Equal(t, p.Data, packet.Data)
//
// 	// Test: drop dublicates
// 	p = newPacket(2)
//
// 	n, err = p.Encode(buf)
// 	require.NoError(t, err)
// 	client.Write(buf[:n])
// 	h.Receive(packets)
//
// 	client.Write(buf[:n])
// 	h.Receive(packets)
//
// 	time.Sleep(100 * time.Millisecond)
// 	require.Equal(t, 1, len(packets))
//
// 	// Test: check seq ack ackbits from many client
// }
