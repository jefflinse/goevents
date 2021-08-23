package goevents

import (
	"time"
)

type Event interface {
	UID() string
	At() *time.Time
	Data() ([]byte, error)
}

type EventHandler func(event Event) error

type DefaultEvent struct {
	Name         string
	DispatchedAt *time.Time
	EventData    []byte
}

func (e *DefaultEvent) UID() string {
	return e.Name
}

func (e *DefaultEvent) At() *time.Time {
	return e.DispatchedAt
}

func (e *DefaultEvent) Data() ([]byte, error) {
	return e.EventData, nil
}

var _ Event = &DefaultEvent{}
