package packet

import "fmt"

func (pubH *PublicHeader) String() string {
	return fmt.Sprintf("Connection id: %d Seq: %d ", pubH.ConnectionID, pubH.Sequence)
}

func (privH *PrivateHeader) String() string {
	return fmt.Sprintf("Ack: %d AckBits: %b Type: %s ", privH.Ack, privH.Ack, privH.Type.String())
}

func (p *Packet) String() string {
	return p.PublicHeader.String() + p.PrivateHeader.String() + fmt.Sprintf("DATA: %q", p.Data)
}
