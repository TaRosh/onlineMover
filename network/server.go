package network

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/TaRosh/online_mover/game"
	"github.com/TaRosh/online_mover/udp"
)

const TickRate = 1

type server struct {
	transport udp.Host
	// maybe ID -> addr -> get to transport layer
	// and addr to Id -> to game layer
	addrToID map[string]game.PlayerID
	idToAddr map[game.PlayerID]*net.UDPAddr
	// next connection ID which will be player Id
	nextConnectionID game.PlayerID
	sentBuf          []byte

	mu sync.RWMutex
}

type NetworkServer interface {
	SendSnapshot(id game.PlayerID, data []byte) error
	// SendPlayerID(id game.PlayerID, data []byte) error
	Receive(inputs chan<- game.Input, connectionEvents chan<- game.PlayerID)
}

func (s *server) isAddrExist(addr *net.UDPAddr) bool {
	s.mu.RLock()
	_, exist := s.addrToID[addr.String()]
	defer s.mu.RUnlock()
	return exist
}

func (s *server) isIDExist(id game.PlayerID) bool {
	s.mu.RLock()
	_, exist := s.idToAddr[id]
	defer s.mu.RUnlock()
	return exist
}

// set two maps and increment id counter
func (s *server) newConnection(addr *net.UDPAddr) game.PlayerID {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id := s.nextConnectionID
	s.addrToID[addr.String()] = id
	s.idToAddr[id] = addr
	s.nextConnectionID += 1
	return id
}

func (s *server) processPacket(packet *udp.Packet, inputHere chan<- game.Input, connectionEventHere chan<- game.PlayerID) {
	// decide which type
	// and what to do next
	// case 1 request for id
	// case 2 player input
	// default drop
	switch packet.Type {
	case udp.PacketInvalid:
		return

	// event for game player connect chan
	case udp.PacketConnect:
		exist := s.isAddrExist(packet.Addr)
		if exist {
			// TODO: what do if old player send connect
			// for now drop
			return
		}
		id := s.newConnection(packet.Addr)
		connectionEventHere <- id
		playerPacketWithID := PlayerIDPacket{
			ID: uint32(id),
		}
		n, err := playerPacketWithID.Encode(s.sentBuf)
		if err != nil || n == 0 {
			// TODO: think about it
			panic(err)
		}

		err = s.sentPlayerID(packet.Addr, s.sentBuf[:n])
		// TODO: think what should do?
		// maybe resend becouse player don't get his id
		// but we allready add him to our game
		// some backup plan to remove if can't sand id
		if err != nil {
			panic(err)
		}
	case udp.PacketInput:

		// for inputs -> game -> player inputs chan
		var input game.Input
		err := input.Decode(packet.Data)
		if err != nil {
			// corruption data in packet
			// just skip
			return
		}
		inputHere <- input

		// TODO: think about disconect user

	}
}

func (s *server) sentPlayerID(addr *net.UDPAddr, data []byte) error {
	err := s.transport.Sent(addr, udp.PacketAccept, data)
	return err
}

// i think for now will use chan for each event from client
func (s *server) Receive(inputHere chan<- game.Input, connectEventHere chan<- game.PlayerID) {
	sendPacketHere := make(chan udp.Packet, 1024)
	go func() {
		for {
			err := s.transport.Receive(sendPacketHere)
			if err != nil {
				log.Println(err)
			}
		}
	}()
	for packet := range sendPacketHere {
		s.processPacket(&packet, inputHere, connectEventHere)
	}
}

func (s *server) SendSnapshot(id game.PlayerID, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	exist := s.isIDExist(id)
	if !exist {
		return fmt.Errorf("Network:SendSnapsho: no addr for this id: %d", id)
	}
	s.mu.RLock()
	addr := s.idToAddr[id]
	s.mu.RUnlock()
	err := s.transport.Sent(addr, udp.PacketSnapshot, data)

	return err
}

var _ NetworkServer = new(server)

func NewServer(host, port string) (*server, error) {
	serv := server{
		addrToID:         make(map[string]game.PlayerID),
		idToAddr:         make(map[game.PlayerID]*net.UDPAddr),
		nextConnectionID: 0,
		sentBuf:          make([]byte, 1024),
	}
	transport, err := udp.NewServer(host, port)
	if err != nil {
		return nil, err
	}
	serv.transport = transport
	serv.sentBuf = make([]byte, 1024)

	return &serv, nil
}
