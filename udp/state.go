package udp

type state struct {
	id             uint32
	lastIDReceived uint32
	packetsIGot    uint32
}
