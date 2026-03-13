package network

import (
	"log"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/udp"
)

type client struct {
	transport udp.Host
}

type NetworkClient interface {
	SendInput(data []byte) error
	SendPlayerConnectionRequest() error
	Receive(snapshots chan<- game.Snapshot, connectionEvent chan<- game.PlayerID)
}

func (c *client) SendPlayerConnectionRequest() error {
	err := c.transport.Sent(c.transport.GetAddr(), udp.PacketConnect, nil)
	return err
}

func (c *client) SendInput(data []byte) error {
	err := c.transport.Sent(c.transport.GetAddr(), udp.PacketInput, data)
	return err
}

func (c *client) processPacket(packet *udp.Packet, snapshots chan<- game.Snapshot, conectionEvent chan<- game.PlayerID) {
	switch packet.Type {
	case udp.PacketInvalid:
		return
	case udp.PacketAccept:
		var pState PlayerIDPacket
		err := pState.Decode(packet.Data)
		if err != nil {
			// TODO: i think just return = corrupted data
			// but for dev panic is ok
			panic(err)
		}
		conectionEvent <- game.PlayerID(pState.ID)
	case udp.PacketSnapshot:
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

func (c *client) Receive(snapshots chan<- game.Snapshot, connectionEvent chan<- game.PlayerID) {
	var sendPacketHere chan udp.Packet
	go func() {
		for {
			err := c.transport.Receive(sendPacketHere)
			// TODO: think about case when error from receive
			log.Println(err)
			continue
		}
	}()
	for packet := range sendPacketHere {
		c.processPacket(&packet, snapshots, connectionEvent)
	}
}

func NewClient(host, port string) (*client, error) {
	c := client{}
	transport, err := udp.NewClient(host, port)
	if err != nil {
		return nil, err
	}
	c.transport = transport
	return &c, nil
}
