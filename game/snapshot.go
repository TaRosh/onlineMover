package game

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/quasilyte/gmath"
)

// PlayerState size is 36
type PlayerState struct {
	ID       uint32
	Position gmath.Vec
	Velocity gmath.Vec
}

// Snapshot size for len of slice uint16
type Snapshot struct {
	Tick    uint32
	Players []PlayerState
}

func (s *Snapshot) Encode(buf []byte) (int, error) {
	// tick -> len(players) -> player state
	binary.BigEndian.PutUint32(buf[0:4], s.Tick)
	// uint16 for len(players)
	binary.BigEndian.PutUint16(buf[4:6], uint16(len(s.Players)))
	// step is 36 byte
	offset := 6
	for _, playerState := range s.Players {
		binary.BigEndian.PutUint32(buf[offset:], playerState.ID)
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
	s.Tick = binary.BigEndian.Uint32(data[0:4])
	counter := binary.BigEndian.Uint16(data[4:])
	offset := 6
	fmt.Println("Counter", counter)
	s.Players = make([]PlayerState, counter)
	for i := range counter {
		p := PlayerState{}
		p.ID = binary.BigEndian.Uint32(data[offset:])
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
