package network

import (
	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/game/entities"
	p "github.com/TaRosh/online_mover/udp/packet"
)

func (c *client) processPacket(packet *p.Packet, snapshots chan<- game.Snapshot, events chan<- game.Event) {
	switch packet.Type {
	case p.Invalid:
		return
	case p.Accept:
		var resp PlayerConnect
		err := resp.Decode(packet.Data)
		if err != nil {
			// TODO: i think just return = corrupted data
			// but for dev panic is ok
			panic(err)
		}
		// TODO: add world size
		events <- game.EventInitConnection{
			ID:          entities.PlayerID(resp.ID),
			WorldWidth:  resp.WorldWidth,
			WorldHeight: resp.WorldHeight,
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
