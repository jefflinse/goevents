package goevents

type Event interface {
	Data() ([]byte, error)
}

type EventHandler func(event Event) error
