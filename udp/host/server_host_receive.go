package host

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	inner "github.com/TaRosh/online_mover/udp/connection"
	"github.com/TaRosh/online_mover/udp/packet"
)

var ErrNoConnection = errors.New("can't get connection by id")

// create random id for connection
func (sh *serverHost) allocateID() uint32 {
	var b [4]byte
	for {
		rand.Read(b[:])
		id := binary.BigEndian.Uint32(b[:])
		if id != 0 {
			_, exist := sh.connections[id]
			if !exist {
				return id
			}
		}
	}
}

func (sh *serverHost) getConnection(connID uint32, playerAddr *net.UDPAddr) *inner.Conn {
	var conn *inner.Conn
	sh.mu.Lock()
	conn = sh.connections[connID]
	// if true ( no conn and id in packet = 0) -> new connection
	if conn == nil && connID == 0 {
		conn = inner.New(playerAddr)
		id := sh.allocateID()
		conn.SetID(id)
		sh.connections[id] = conn
	}
	sh.mu.Unlock()
	return conn
}

func (sh *serverHost) Receive(sendPacketHere chan<- packet.Packet) error {
	buf := sh.getBuf()
	defer sh.putBuf(buf)

	n, playerAddr, err := sh.conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}

	// 1. Get connection
	// Try get id from public header field connection ID
	var pubHeader packet.PublicHeader
	_, err = pubHeader.Decode(buf[:n])
	// without public header can't recevie -> return
	// dropp packet
	if err != nil {
		return fmt.Errorf("serverHost:Receive: %w", err)
	}
	var conn *inner.Conn

	// try get connection by id
	connID := pubHeader.ConnectionID
	conn = sh.getConnection(connID, playerAddr)

	// if conn still nill then:
	// packet have connection ID
	// but no such connection on server
	// maybe deleted by some bug
	// or someone try change it
	if conn == nil {
		return ErrNoConnection
	}

	// we change addr in host.receive
	// by middleware change Addr
	packet := sh.getPacket()
	defer sh.putPacket(packet)

	packet.Addr = playerAddr

	err = sh.host.receive(conn, buf[:n], packet)
	if err != nil {
		return err
	}
	if packet != nil {
		sendPacketHere <- *packet
	}
	return nil
}
