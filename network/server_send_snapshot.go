package network

import (
	"fmt"

	"github.com/TaRosh/online_mover/game/entities"
	p "github.com/TaRosh/online_mover/udp/packet"
)

func (s *server) SendSnapshot(id entities.PlayerID, data []byte) error {
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
