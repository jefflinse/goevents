package goevents

import (
	"time"
)

type Event interface {
	Name() string
	Data() ([]byte, error)
}

type EventHandler func(event Event) error

type DefaultEvent struct {
	EventName    string
	DispatchedAt *time.Time
	EventData    []byte
}

func (e *DefaultEvent) Name() string {
	return e.EventName
}

func (e *DefaultEvent) Data() ([]byte, error) {
	return e.EventData, nil
}

var _ Event = &DefaultEvent{}
