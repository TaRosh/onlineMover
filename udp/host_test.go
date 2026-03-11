package udp

import (
	"testing"
)

func TestSendReceive(t *testing.T) {
}

func TestAckBits(t *testing.T) {
	// test ackbits when newer packet arrive
	// h, _ := NewServer("3000")
	// h.state.lastIDReceived = 100
	// h.state.packetsIGot = 0b0111
	// p := Packet{
	// 	Header: Header{
	// 		Sequence: 101,
	// 	},
	// }
	// h.processPacket(&p)
	// if h.state.packetsIGot != 0b1111 {
	// 	t.Errorf("want %b, got %b", 0b1111, h.packetsIGot)
	// }
	// if h.state.lastIDReceived != 101 {
	// 	t.Errorf("not change last received id to newer packet")
	// }
	// // test ackbits when old packet arrive
	// h.state.lastIDReceived = 100
	// h.state.packetsIGot = 0b1011
	// p.Header.Sequence = 98
	// h.processPacket(&p)
	// if h.state.packetsIGot != 0b1111 {
	// 	t.Errorf("want %b, got %b", 0b1111, h.packetsIGot)
	// }
	// if h.state.lastIDReceived != 100 {
	// 	t.Errorf("change last received id when old packet arrive")
	// }
}
