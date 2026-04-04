package host

import "github.com/TaRosh/online_mover/udp/packet"

func (h *host) getPacket() *packet.Packet {
	return h.packetPool.Get().(*packet.Packet)
}

func (h *host) putPacket(p *packet.Packet) {
	*p = packet.Packet{}
	h.packetPool.Put(p)
}

func (h *host) getBuf() []byte {
	return h.bufferPool.Get().([]byte)
}

func (h *host) putBuf(buf []byte) {
	buf = buf[:cap(buf)]
	h.bufferPool.Put(buf)
}
