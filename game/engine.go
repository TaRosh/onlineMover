package game

import "github.com/quasilyte/gmath"

const speed = 2

func ApplyInput(p *Player, input Input) {
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
