package network

import "github.com/TaRosh/online_mover/game"

type NetworkClient interface {
	SendInput(data []byte) error
	SendPlayerConnectionRequest() error
	// Connect() error
	Receive(snapshots chan<- game.Snapshot, events chan<- game.Event)
	// when no new packet from these players in 5 sec.
	// notify
	// delete player connection from transport layer
	// CleanUp(playerID game.PlayerID)
}
type NetworkServer interface {
	SendSnapshot(id game.PlayerID, data []byte) error
	DeletePlayer(id game.PlayerID)
	CheckTimeouts(events chan<- game.Event)
	Receive(inputs chan<- game.Input, events chan<- game.Event)
}
