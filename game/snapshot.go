package game

import (
	"encoding/binary"
)

// PlayerState size is 36
type PlayerState struct {
	ID PlayerID
	X  int32
	Y  int32
	VX int16
	VY int16
	// Position gmath.Vec
	// Velocity gmath.Vec
}

type ProjectileState struct {
	ID uint32
	X  int32
	Y  int32
	VX int16
	VY int16
	// Position gmath.Vec
	// Velocity gmath.Vec
}

// type WorldStat struct {
// 	BaseX int32
// 	BaseY int32
// }

// Snapshot size for len of slice uint16
type Snapshot struct {
	Tick          uint32
	LastInputTick uint32
	Players       []PlayerState
	Projectiles   []ProjectileState
}

func (s *Snapshot) Encode(buf []byte) (int, error) {
	var offset int
	// tick -> len(players) -> player state
	binary.BigEndian.PutUint32(buf[offset:], s.Tick)
	offset += 4
	binary.BigEndian.PutUint32(buf[offset:], s.LastInputTick)
	offset += 4

	// uint16 for len(players)
	binary.BigEndian.PutUint16(buf[offset:], uint16(len(s.Players)))
	// step is 36 byte
	offset += 2
	for _, playerState := range s.Players {
		binary.BigEndian.PutUint32(buf[offset:], uint32(playerState.ID))
		offset += 4
		binary.BigEndian.PutUint32(buf[offset:], uint32(playerState.X))
		offset += 4
		binary.BigEndian.PutUint32(buf[offset:], uint32(playerState.Y))
		offset += 4
		binary.BigEndian.PutUint16(buf[offset:], uint16(playerState.VX))
		offset += 2
		binary.BigEndian.PutUint16(buf[offset:], uint16(playerState.VY))
		offset += 2
	}
	// uint16 for len(players)
	binary.BigEndian.PutUint16(buf[offset:], uint16(len(s.Projectiles)))
	// step is 36 byte
	offset += 2
	for _, projectileState := range s.Projectiles {
		binary.BigEndian.PutUint32(buf[offset:], projectileState.ID)
		offset += 4
		binary.BigEndian.PutUint32(buf[offset:], uint32(projectileState.X))
		offset += 4
		binary.BigEndian.PutUint32(buf[offset:], uint32(projectileState.Y))
		offset += 4
		binary.BigEndian.PutUint16(buf[offset:], uint16(projectileState.VX))
		offset += 2
		binary.BigEndian.PutUint16(buf[offset:], uint16(projectileState.VY))
		offset += 2
	}

	return offset, nil
}

func (s *Snapshot) Decode(data []byte) error {
	var offset int
	s.Tick = binary.BigEndian.Uint32(data[offset:])
	offset += 4
	s.LastInputTick = binary.BigEndian.Uint32(data[offset:])
	offset += 4
	counter := binary.BigEndian.Uint16(data[offset:])
	offset += 2
	s.Players = make([]PlayerState, counter)
	for i := range counter {
		p := PlayerState{}
		p.ID = PlayerID(binary.BigEndian.Uint32(data[offset:]))
		offset += 4

		p.X = int32(binary.BigEndian.Uint32(data[offset:]))
		offset += 4
		p.Y = int32(binary.BigEndian.Uint32(data[offset:]))
		offset += 4
		p.VX = int16(binary.BigEndian.Uint16(data[offset:]))
		offset += 2
		p.VY = int16(binary.BigEndian.Uint16(data[offset:]))
		offset += 2

		s.Players[i] = p
	}
	counter = binary.BigEndian.Uint16(data[offset:])
	offset += 2
	s.Projectiles = make([]ProjectileState, counter)
	for i := range counter {
		p := ProjectileState{}
		p.ID = binary.BigEndian.Uint32(data[offset:])
		offset += 4

		p.X = int32(binary.BigEndian.Uint32(data[offset:]))
		offset += 4
		p.Y = int32(binary.BigEndian.Uint32(data[offset:]))
		offset += 4
		p.VX = int16(binary.BigEndian.Uint16(data[offset:]))
		offset += 2
		p.VY = int16(binary.BigEndian.Uint16(data[offset:]))
		offset += 2

		s.Projectiles[i] = p
	}
	return nil
}
