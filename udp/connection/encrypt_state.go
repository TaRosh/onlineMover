package connection

type SecureState int

const (
	Unsecure SecureState = iota
	Wait
	Secure
)

func (s SecureState) String() string {
	switch s {
	case Unsecure:
		return "Unsecure"
	case Wait:
		return "Wait"
	case Secure:
		return "Secure"
	}
	return "INVALID"
}

// state about secure or unsecure connection
func (conn *Conn) GetEncryptionState() SecureState {
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	return conn.state
}

func (conn *Conn) SetEncryptionState(state SecureState) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.state = state
}
