package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnapshotEncodeDecode(t *testing.T) {
	original := Snapshot{
		Tick:          4,
		LastInputTick: 2,
		Players: []PlayerState{
			{
				ID: 0,
				X:  1,
				Y:  1,
				VX: 1,
				VY: 0,
			},
			{
				ID: 1,
				X:  0,
				Y:  0,
				VX: 0,
				VY: 1,
			},
		},
		Projectiles: []ProjectileState{
			{
				ID: 0,
				X:  1,
				Y:  0,
				VX: 0,
				VY: 1,
			},
			{
				ID: 1,
				X:  0,
				Y:  1,
				VX: 1,
				VY: 0,
			},
		},
	}
	buf := make([]byte, 1024)
	n, err := original.Encode(buf)
	require.NoError(t, err)
	require.NotEqual(t, 0, n)
	decoded := Snapshot{}

	err = decoded.Decode(buf[:n])
	require.NoError(t, err)
	require.Equal(t, original, decoded)
}
