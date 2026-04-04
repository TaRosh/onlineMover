package network

import (
	"sync"

	"github.com/TaRosh/online_mover/udp/host"
	"github.com/TaRosh/online_mover/game"
)

const TickRate = 1

type server struct {
	transport host.ServerHost
	// maybe player ID -> connection id -> get to transport layer
	// and connection id to  palyer id -> to game layer
	connToPlayer map[uint32]game.PlayerID
	playerToConn map[game.PlayerID]uint32
	// next connection ID which will be player Id
	nextConnectionID game.PlayerID
	sentBuf          []byte

	mu sync.RWMutex
}

var _ NetworkServer = new(server)

func NewServer(hst, port string, secretKey []byte) (*server, error) {
	serv := server{
		connToPlayer:     map[uint32]game.PlayerID{},
		playerToConn:     map[game.PlayerID]uint32{},
		nextConnectionID: 0,

		sentBuf: make([]byte, 1200),
		mu:      sync.RWMutex{},
	}
	transport, err := host.NewServer(hst, port)
	if err != nil {
		return nil, err
	}
	serv.transport = transport

	return &serv, nil
}

// DeleteConn(id *net.UDPAddr)
func (s *server) DeletePlayer(id game.PlayerID) {
	s.transport.DeleteConn(s.playerToConn[id])
}

func (s *server) CheckTimeouts(events chan<- game.Event) {
	// events -> network
	// network: -> transport
	// transport:addr -> network
	// netwokr:id -> events

	disconnectedConnection := make(chan uint32, 24)

	go s.transport.CheckTimeouts(disconnectedConnection)
	for id := range disconnectedConnection {
		events <- game.Event{
			Type: game.EventNoAnswerFromClient,

			ID: s.connToPlayer[id],
		}
	}
}

func (s *server) isPlayerExist(id uint32) bool {
	s.mu.RLock()
	_, exist := s.connToPlayer[id]
	defer s.mu.RUnlock()
	return exist
}

func (s *server) isConnectionExist(id game.PlayerID) bool {
	s.mu.RLock()
	_, exist := s.playerToConn[id]
	defer s.mu.RUnlock()
	return exist
}

// set two maps and increment id counter
func (s *server) newConnection(connID uint32) game.PlayerID {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id := s.nextConnectionID
	s.connToPlayer[connID] = id
	s.playerToConn[id] = connID
	s.nextConnectionID += 1
	return id
}
