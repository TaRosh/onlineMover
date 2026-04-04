package host

import (
	"net"
	"time"

	inner "github.com/TaRosh/online_mover/udp/connection"
	"github.com/TaRosh/online_mover/udp/packet"
)

type clientHost struct {
	*host
	connection                *inner.Conn
	lastEncryptionRequestSend time.Time
}

func NewClient(h, port string) (*clientHost, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(h, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	clntHost := &clientHost{
		host:       newHost(conn),
		connection: inner.New(nil),
	}
	// set handler for encryption answer from server
	clntHost.host.handleEncrypt[packet.KeyExchangeAnswer] = clntHost.handleEncryptionResponse
	return clntHost, nil
}
