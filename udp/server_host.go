package udp

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type serverHost struct {
	*host
	connections map[uint32]*connection
	addrToConn  map[string]*connection
	mu          sync.RWMutex
}

func (sh *serverHost) Send(id uint32, t packetType, data []byte) error {
	conn, exist := sh.connections[id]
	if !exist {
		return errors.New(fmt.Sprintf("serverHost:Send: no connection with id: %d", id))
	}
	err := sh.send(conn, t, data)
	return err
}

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

func (sh *serverHost) Receive(sendPacketHere chan<- Packet) error {
	buf := sh.bufferPool.Get().([]byte)
	defer func() {
		buf = buf[:cap(buf)]
		sh.bufferPool.Put(buf)
	}()
	n, playerAddr, err := sh.conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}

	// Try get connection via player addr
	key := playerAddr.String()
	sh.mu.RLock()
	conn := sh.addrToConn[key]
	sh.mu.RUnlock()

	if conn != nil {
		return sh.processRawPacket(conn, buf[:n], sendPacketHere)
	}

	// Try get connection via connection id from packet
	var pubHeader PublicHeader
	pubHeader.Decode(buf[:n])
	if pubHeader.ConnectionID != 0 {
		return sh.processRawPacket(conn, buf[:n], sendPacketHere)
	}

	// if can't get connection via address
	// and via connection ID
	// it must be new connection
	if conn == nil {
		fmt.Println("Create new connection")
		conn = newConnection(playerAddr)

		sh.mu.Lock()
		id := sh.allocateID()

		sh.addrToConn[key] = conn
		sh.connections[id] = conn
		conn.id = id

		sh.mu.Unlock()
	}
	err = sh.processRawPacket(conn, buf[:n], sendPacketHere)
	if err != nil {
		return err
	}
	return nil
}

/*
// Give existed or create new connection
func (h *host) getConnection(id uint32, addr *net.UDPAddr) *connection {
	h.mu.RLock()
	conn, exist := h.connections[key]
	h.mu.RUnlock()
	if exist {
		return conn
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	conn, exist = h.connections[key]
	if exist {
		return conn
	}

	conn = newConnection(addr)
	id := h.allocateIDLocked()
	h.connections[key] = conn

	return conn
}
*/

func (sh *serverHost) DeleteConn(id uint32) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	delete(sh.connections, id)
}

func (sh *serverHost) CheckTimeouts(noAnswerAddrHere chan<- *net.UDPAddr) {
	defer close(noAnswerAddrHere)
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	for _, conn := range sh.connections {
		conn.mutex.RLock()
		lastTimeHeared := conn.lastHeared
		fmt.Println(time.Since(lastTimeHeared))
		conn.mutex.RUnlock()
		if time.Since(lastTimeHeared) > 5*time.Second {
			if conn.Addr == nil {
				panic("checkTimeouts")
			}
			noAnswerAddrHere <- conn.Addr
		}
	}
}

func NewServer(h, port string) (*serverHost, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(h, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	/*
		connections:
		make(map[string]*connection),
	*/
	srvHost := &serverHost{
		host:        newHost(conn),
		connections: make(map[uint32]*connection),
		addrToConn:  make(map[string]*connection),
	}
	return srvHost, nil
}
