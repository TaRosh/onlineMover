package main

func stateEvents(w *World) stateFn {
	for {
		select {
		case event := <-w.incomingEvents:
			w.processEvent(event)
		default:
			return stateInputs
		}
	}
}
