package game

import "github.com/TaRosh/online_mover/game/entities"

type EventInitConnection struct {
	ID          entities.PlayerID
	WorldWidth  uint16
	WorldHeight uint16
}

func (e EventInitConnection) isEvent() {}

// type Event struct {
// 	Type EventType
// 	ID   PlayerID
// }
