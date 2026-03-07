package udp

import (
	"fmt"
	"log"
	"net"
	"time"
)

type client struct {
	packetsSend map[uint32]SentPacket
	conn        *net.UDPConn
	state
	smoothedRTT time.Duration
	lastCleanup time.Time
	buf         []byte
}

type NetworkClient interface {
	SendInput(data []byte) error
}

func (c *client) SendInput(data []byte) error {
	err := c.Write(inputPacket, data)
	return err
}

func (c *client) Write(t packetType, data []byte) error {
	packet := NewPacket(c.state.id, c.state.lastIDReceived, c.state.packetsIGot, t, data)
	n, err := packet.Encode(c.buf)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(c.buf[:n])
	if err != nil {
		return err
	}
	c.packetsSend[c.state.id] = SentPacket{
		sendedWhen: time.Now(),
		delivered:  false,
	}
	// fmt.Printf("sending Seq: %d Ack: %d\n", packet.Sequence, packet.Ack)
	c.id += 1
	return nil
}

// check is packet id is old or new
func isNewer(a, b uint32) bool {
	return int32(a-b) > 0
}

func (c *client) calcRtt(sample time.Duration) {
	if c.smoothedRTT == 0 {
		c.smoothedRTT = sample
	} else {
		c.smoothedRTT = time.Duration(float64(c.smoothedRTT)*0.9 + float64(sample)*0.1)
	}
}

// do ack check and clean up in one loop
// may create separate function and add mutex
func (c *client) processAck(header Header) {
	// check ack & ack bits for get rtt and calculate average rtt
	if sentPacket, ok := c.packetsSend[header.Ack]; ok {
		if !sentPacket.delivered {
			rawRtt := time.Since(sentPacket.sendedWhen)
			c.calcRtt(rawRtt)
			delete(c.packetsSend, header.Ack)
		}
	}
	for i := range uint32(32) {
		if (header.AckBits>>i)&1 == 1 {
			if sentPacket, ok := c.packetsSend[header.Ack-i]; ok {
				if !sentPacket.delivered {
					rawRtt := time.Since(sentPacket.sendedWhen)
					c.calcRtt(rawRtt)
					delete(c.packetsSend, header.Ack)
				}
			}
		}
	}
	// clean up the map for sended packets
	// delete data about packet if  not arrived after 3 second
	// launch every second
	if time.Since(c.lastCleanup) > time.Second {
		for id, sentPacket := range c.packetsSend {
			if time.Since(sentPacket.sendedWhen) > time.Second*3 {
				delete(c.packetsSend, id)
			}
		}
		c.lastCleanup = time.Now()
	}
}

func (c *client) processPacket(p *Packet) {
	c.processAck(p.Header)
	// if packet id is new ( < then last id i got)
	// check is inside our check table ( uint32 ) or we lost
	// 32 packet and it is newer
	// do new check table or shift for it's diff ( currentId - lastId )
	if isNewer(p.Sequence, c.lastIDReceived) {
		shift := p.Sequence - c.lastIDReceived
		if shift >= 32 {
			c.packetsIGot = 0
		} else {
			c.packetsIGot <<= shift
		}
		c.packetsIGot |= 1
		c.lastIDReceived = p.Sequence

	} else {
		diff := c.lastIDReceived - p.Sequence
		if diff < 32 && diff >= 1 {
			bitIndex := diff - 1
			c.packetsIGot |= (1 << bitIndex)
		}
		// if diff >= 32 {
		// 	// TODO: package to old what to do?
		// }
	}
	// check is packet in check table
}

func (c *client) Read() {
	for {

		n, _, err := c.conn.ReadFromUDP(c.buf)
		if err != nil {
			log.Println(err)
		}
		packet := Packet{}
		err = packet.Decode(c.buf[:n])
		if err != nil {
			log.Println("ERROR:", err)
		}

		fmt.Printf("Receive: Seq: %d, Ack: %d, AckBits:%d\n", packet.Sequence, packet.Ack, packet.AckBits)
		c.processPacket(&packet)
		// fmt.Printf("Last come %d\n", c.lastIDReceived)
		// for i := range 32 {
		// 	fmt.Printf("[SEND] packet %d -> %b\n", packet.Ack-uint32(i), byte(packet.AckBits>>i&0x01))
		// }

	}
}

func NewClient(host, port string) (*client, error) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	c := client{
		conn: conn,
		// seqNum:    0xaabbccdd,
		state: state{
			id: 1,
		},
		buf:         make([]byte, 2048),
		packetsSend: make(map[uint32]SentPacket),
	}

	return &c, nil
}
