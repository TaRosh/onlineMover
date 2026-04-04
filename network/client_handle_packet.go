package network

import (
	p "github.com/TaRosh/online_mover/udp/packet"
	"github.com/TaRosh/online_mover/game"
)

func (c *client) processPacket(packet *p.Packet, snapshots chan<- game.Snapshot, events chan<- game.Event) {
	switch packet.Type {
	case p.Invalid:
		return
	case p.Accept:
		var pState PlayerIDPacket
		err := pState.Decode(packet.Data)
		if err != nil {
			// TODO: i think just return = corrupted data
			// but for dev panic is ok
			panic(err)
		}
		events <- game.Event{
			Type: game.EventConnection,
			ID:   game.PlayerID(pState.ID),
		}
	case p.Snapshot:
		var snapshot game.Snapshot
		err := snapshot.Decode(packet.Data)
		if err != nil {
			// TODO: i think just return = corrupted data
			// but for dev panic is ok
			panic(err)
		}
		snapshots <- snapshot

	}
}
