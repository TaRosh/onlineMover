package game

import (
	"encoding/binary"
	"math"

	"github.com/quasilyte/gmath"
)

// PlayerState size is 36
type PlayerState struct {
	ID       PlayerID
	Position gmath.Vec
	Velocity gmath.Vec
}

// Snapshot size for len of slice uint16
type Snapshot struct {
	Tick          uint32
	LastInputTick uint32
	Players       []PlayerState
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
		binary.BigEndian.PutUint64(buf[offset:], math.Float64bits(playerState.Position.X))
		offset += 8
		binary.BigEndian.PutUint64(buf[offset:], math.Float64bits(playerState.Position.Y))
		offset += 8
		binary.BigEndian.PutUint64(buf[offset:], math.Float64bits(playerState.Velocity.X))
		offset += 8
		binary.BigEndian.PutUint64(buf[offset:], math.Float64bits(playerState.Velocity.Y))
		offset += 8
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

		p.Position.X = math.Float64frombits(binary.BigEndian.Uint64(data[offset:]))
		offset += 8
		p.Position.Y = math.Float64frombits(binary.BigEndian.Uint64(data[offset:]))
		offset += 8
		p.Velocity.X = math.Float64frombits(binary.BigEndian.Uint64(data[offset:]))
		offset += 8
		p.Velocity.X = math.Float64frombits(binary.BigEndian.Uint64(data[offset:]))
		offset += 8

		s.Players[i] = p
	}
	return nil
}
