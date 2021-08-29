package goevents

import (
	"strings"
	"time"
)

type Event interface{}

func EventName(e Event) string {
	return strings.TrimSuffix(typeName(e), "Event")
}

type EventContext struct {
	Type         string
	DispatchedAt time.Time
	Event        Event
}

type EventHandler func(e *EventContext) error

var NoopEventHandler = func(e *EventContext) error {
	return nil
}
