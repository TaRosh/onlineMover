package network

import (
	"fmt"

	p "github.com/TaRosh/online_mover/udp/packet"
	"github.com/TaRosh/online_mover/game"
)

func (s *server) sentPlayerID(connID uint32, data []byte) error {
	err := s.transport.Send(connID, p.Accept, data)
	return err
}

func (s *server) SendSnapshot(id game.PlayerID, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	exist := s.isConnectionExist(id)
	if !exist {
		return fmt.Errorf("Network:SendSnapsho: no connection for this id: %d", id)
	}
	s.mu.RLock()
	addr := s.playerToConn[id]
	s.mu.RUnlock()
	err := s.transport.Send(addr, p.Snapshot, data)

	return err
}
