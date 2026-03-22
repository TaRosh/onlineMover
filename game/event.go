package game

type EventType uint32

const (
	WrongType EventType = iota
	EventConnection
	EventNoAnswerFromClient
)

type Event struct {
	Type EventType
	ID   PlayerID
}
