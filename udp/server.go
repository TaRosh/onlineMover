package udp

import (
	"log"
	"net"
)

type server struct {
	conn *net.UDPConn
	done bool
	buf  []byte
}

func (s *server) Listen() {
	defer s.Close()
	for {
		// TODO: check buf used as ring, not append
		n, clientAddr, err := s.conn.ReadFromUDP(s.buf)
		if err != nil {
			log.Println("server:Listen:read", err)
			continue
		}
		log.Printf("recv %d bytes from %v\n", n, clientAddr)
		// echo back
		_, err = s.conn.WriteToUDP(s.buf[:n], clientAddr)
		if err != nil {
			log.Println("server:Listen:write:", err)
			continue
		}

	}
}

func (s *server) Close() {
	s.conn.Close()
}

func New(port string) (*server, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("", port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &server{
		conn: conn,
		buf:  make([]byte, 2048),
	}, nil
}
