package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/main_client/player"
)

func (g *Game) connectPlayerToServer() {
	tries := 4
	// time after we resend connection request
	resendAfter := time.Tick(time.Second * 2)
	// for range tries {
TRY_GET_PLAYER_ID_AGAIN:
	if tries == 0 {
		panic("could'n get player id from server. I'm done! Shut Done")
	}
	err := g.Network.SendPlayerConnectionRequest()
	// another try when network problem
	if err != nil {
		tries--
		time.Sleep(time.Millisecond * 300)
		goto TRY_GET_PLAYER_ID_AGAIN
	}
	// another try when no answer from server
	for {
		select {
		case e := <-g.events:

			switch event := e.(type) {
			case game.EventInitConnection:
				id := event.ID
				fmt.Println("MY ID", id)
				// add player
				player, _ := player.New(id, color.White)
				// TODO: check error
				g.localPlayer = player
				g.players[id] = player
				//  set world size
				g.WorldWidth = int(event.WorldWidth)
				g.WorldHeight = int(event.WorldHeight)
				// g.world = ebiten.NewImage(g.WorldWidth, g.WorldHeight)
			}
			return
		case <-resendAfter:
			tries--
			goto TRY_GET_PLAYER_ID_AGAIN
		}
	}
}
