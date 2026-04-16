package game

import "github.com/TaRosh/online_mover/game/entities"

type EventNoAnswerFromClient struct {
	ID entities.PlayerID
}

func (e EventNoAnswerFromClient) isEvent() {}
