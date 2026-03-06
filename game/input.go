package game

import "encoding/binary"

// size 4 + 4
type Input struct {
	Tick  uint32
	Up    bool
	Down  bool
	Left  bool
	Right bool
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func (i *Input) Encode(buf []byte) (int, error) {
	binary.BigEndian.AppendUint32(buf[0:4], i.Tick)
	buf[4] = boolToByte(i.Up)
	buf[5] = boolToByte(i.Down)
	buf[6] = boolToByte(i.Left)
	buf[7] = boolToByte(i.Right)

	return 8, nil
}

func (i *Input) Decode(data []byte) error {
	i.Tick = binary.BigEndian.Uint32(data[0:4])
	i.Up = data[4] != 0
	i.Down = data[5] != 0
	i.Left = data[6] != 0
	i.Right = data[7] != 0
	return nil
}
