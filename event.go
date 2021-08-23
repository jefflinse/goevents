package goevents

import (
	"encoding/json"
)

type Event interface {
	Data() ([]byte, error)
}

type EventHandler func(event Event) error

// A JSONEvent is an event whose data is simply the event object marshalled as JSON.
type JSONEvent struct{}

func (e JSONEvent) Data() ([]byte, error) {
	return json.Marshal(e)
}

var _ Event = &JSONEvent{}
