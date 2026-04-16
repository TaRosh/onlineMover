package game

import (
	"github.com/TaRosh/online_mover/game/entities"
	"github.com/quasilyte/gmath"
)

const speed = 5

func ApplyInput(p *entities.Player, input Input) {
	var force gmath.Vec
	if input.Buttons&InputLeft != 0 {
		force.X = -speed
	}
	if input.Buttons&InputRight != 0 {
		force.X = speed
	}
	if input.Buttons&InputDown != 0 {
		force.Y = speed
	}
	if input.Buttons&InputUp != 0 {
		force.Y = -speed
	}
	// p.ApplyForce(force)
	p.Velocity = force
}
