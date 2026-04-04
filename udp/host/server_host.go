package host

import (
	"net"
	"sync"
	"time"

	inner "github.com/TaRosh/online_mover/udp/connection"
	"github.com/TaRosh/online_mover/udp/packet"
)

type serverHost struct {
	*host
	connections map[uint32]*inner.Conn
	addrToConn  map[string]*inner.Conn
	mu          sync.RWMutex
}

func (sh *serverHost) DeleteConn(id uint32) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	delete(sh.connections, id)
}

func (sh *serverHost) CheckTimeouts(noAnswerIDHere chan<- uint32) {
	defer close(noAnswerIDHere)
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	for _, conn := range sh.connections {
		lastTimeHeared := conn.GetLastTimeActive()
		// fmt.Println(time.Since(lastTimeHeared))
		if time.Since(lastTimeHeared) > 5*time.Second {
			if conn.Addr == nil {
				panic("checkTimeouts")
			}
			noAnswerIDHere <- conn.GetID()
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
	srvHost := &serverHost{
		host:        newHost(conn),
		connections: make(map[uint32]*inner.Conn),
		addrToConn:  make(map[string]*inner.Conn),
	}
	// set handler for encryption request from client
	srvHost.host.handleEncrypt[packet.KeyExchangeRequest] = srvHost.handleEncryptionRequest
	return srvHost, nil
}
