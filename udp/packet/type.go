package packet

type Type uint8

const (
	Invalid Type = iota
	Input
	Snapshot
	Connect
	KeyExchangeRequest
	KeyExchangeAnswer
	Accept
)

func (p Type) String() string {
	switch p {
	case Invalid:
		return "Invalid"
	case Input:
		return "Input"
	case Snapshot:
		return "Snapshot"
	case Connect:
		return "Connect"
	case KeyExchangeRequest:
		return "KeyExchangeRequest"
	case KeyExchangeAnswer:
		return "KeyExchangeAnswer"
	case Accept:
		return "Accept"
	}
	return ""
}
