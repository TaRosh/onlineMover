package host

import (
	"net"
	"sync"

	inner "github.com/TaRosh/online_mover/udp/connection"
	"github.com/TaRosh/online_mover/udp/packet"
	p "github.com/TaRosh/online_mover/udp/packet"
)

const (
	bufferSize = 1200
)

type host struct {
	conn *net.UDPConn

	bufferPool    sync.Pool
	packetPool    sync.Pool
	handleEncrypt map[packet.Type]func(*inner.Conn, *p.Packet) error
	// decode packet by connection secure state ( unsecure | secure )
	handleConn map[inner.SecureState]func(*inner.Conn, *p.Packet, []byte) error
	middleware []func(*inner.Conn, *p.Packet) error
}

func (h *host) Close() {
	h.conn.Close()
}

func newHost(conn *net.UDPConn) *host {
	h := host{
		conn:          conn,
		handleConn:    make(map[inner.SecureState]func(*inner.Conn, *p.Packet, []byte) error),
		handleEncrypt: make(map[p.Type]func(*inner.Conn, *p.Packet) error),
		middleware:    make([]func(*inner.Conn, *p.Packet) error, 0),
		packetPool: sync.Pool{
			New: func() any {
				return &p.Packet{}
			},
		},
		bufferPool: sync.Pool{
			New: func() any {
				return make([]byte, bufferSize)
			},
		},
	}
	// set handlers
	h.handleConn[inner.Unsecure] = handleUnsecure
	h.handleConn[inner.Secure] = handleSecure
	h.handleConn[inner.Wait] = handleUnsecure
	h.middleware = append(h.middleware, middlewareChangeAddr)
	// encrypt handlers set client/server host
	return &h
}
