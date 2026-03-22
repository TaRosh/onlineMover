package network

import (
	"fmt"
	"log"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/udp"
)

type client struct {
	transport udp.ClientHost
}

type NetworkClient interface {
	SendInput(data []byte) error
	SendPlayerConnectionRequest() error
	Receive(snapshots chan<- game.Snapshot, events chan<- game.Event)
	// when no new packet from these players in 5 sec.
	// notify
	// delete player connection from transport layer
	// CleanUp(playerID game.PlayerID)
}

func (c *client) SendPlayerConnectionRequest() error {
	err := c.transport.Send(udp.PacketConnect, nil)
	return err
}

func (c *client) SendInput(data []byte) error {
	err := c.transport.Send(udp.PacketInput, data)
	return err
}

func (c *client) processPacket(packet *udp.Packet, snapshots chan<- game.Snapshot, events chan<- game.Event) {
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
		events <- game.Event{
			Type: game.EventConnection,
			ID:   game.PlayerID(pState.ID),
		}
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

func (c *client) Receive(snapshots chan<- game.Snapshot, events chan<- game.Event) {
	sendPacketHere := make(chan udp.Packet, 1024)
	go func() {
		for {
			err := c.transport.Receive(sendPacketHere)
			// TODO: think about case when error from receive
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}()
	for packet := range sendPacketHere {
		fmt.Println("RECEIVE packet", packet)
		c.processPacket(&packet, snapshots, events)
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
