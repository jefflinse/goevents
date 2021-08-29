package goevents

import "time"

type Event interface{}

type EventContext struct {
	Type         string
	DispatchedAt time.Time
	Event        Event
}

type EventHandler func(e *EventContext) error

var NoopEventHandler = func(e *EventContext) error {
	return nil
}
