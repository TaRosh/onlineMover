package network

import (
	"github.com/TaRosh/online_mover/game/entities"
	"github.com/TaRosh/online_mover/udp/packet"
)

func (s *server) SendInitConnectionResponse(id entities.PlayerID, w, h int) error {
	req := PlayerConnect{
		ID:          uint32(id),
		WorldWidth:  uint16(w),
		WorldHeight: uint16(h),
	}

	n, err := req.Encode(s.sentBuf)
	if err != nil || n == 0 {
		// TODO: ???
		panic(err)
	}
	connId := s.playerToConn[id]
	err = s.transport.Send(connId, packet.Accept, s.sentBuf[:n])
	// TODO: think what should do?
	// maybe resend becouse player don't get his id
	// but we allready add him to our game
	// some backup plan to remove if can't sand id
	return err
}
