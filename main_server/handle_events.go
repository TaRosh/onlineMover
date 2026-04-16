package main

import (
	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/game/entities"
)

func (w *World) processEvent(e game.Event) {
	switch event := e.(type) {
	case game.EventInitConnection:
		// 1. Create new player
		player := entities.NewPlayer(event.ID, 32, 32)
		w.players[player.ID] = player
		player.Position = w.playerSpawnPoint
		// 2. Send world size
		// and player id
		err := w.network.SendInitConnectionResponse(player.ID, w.Width, w.Height)
		if err != nil {
			// TODO: think what should do?
			// maybe resend becouse player don't get his id
			// but we allready add him to our game
			// some backup plan to remove if can't sand id
			panic(err)
		}

		// 3. notify players about new player
		w.sendSnapshot(game.SnapshotFull)
	case game.EventNoAnswerFromClient:
		w.network.DeletePlayer(event.ID)
		delete(w.players, event.ID)
	}
}
