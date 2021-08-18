package goevents

import (
	"log"
	"time"
)

type Event interface {
	Type() string
	OccurredAt() time.Time
	Data() []byte
}

type EventHandler func(event Event) error

type BasicEvent struct {
	EventType       string
	EventOccurredAt time.Time
	EventData       []byte
}

func (e *BasicEvent) Type() string {
	return e.EventType
}

func (e *BasicEvent) OccurredAt() time.Time {
	return e.EventOccurredAt
}

func (e *BasicEvent) Data() []byte {
	return e.EventData
}

var _ Event = &BasicEvent{}

// An EventBus is anything capable of event pub/sub.
type EventBus interface {
	Publish(event Event) error
	Subscribe(eventType string, handler EventHandler) error
	SubscribeAll(handler EventHandler) error
}

// MemoryEventBus is an in-memory event bus implementation.
type MemoryEventBus struct {
	globalSubscribers []EventHandler
	subscribers       map[string][]EventHandler
}

var _ EventBus = &MemoryEventBus{}

func (bus *MemoryEventBus) Publish(event Event) error {
	log.Printf("[event] %s {%s}\n", event.Type(), string(event.Data()))

	for _, globalHandler := range bus.globalSubscribers {
		if err := globalHandler(event); err != nil {
			return err
		}
	}

	for _, handler := range bus.subscribers[event.Type()] {
		if err := handler(event); err != nil {
			return err
		}
	}

	return nil
}

func (bus *MemoryEventBus) Subscribe(eventType string, handler EventHandler) error {
	if bus.subscribers == nil {
		bus.subscribers = make(map[string][]EventHandler)
	}

	if bus.subscribers[eventType] == nil {
		bus.subscribers[eventType] = make([]EventHandler, 0)
	}

	bus.subscribers[eventType] = append(bus.subscribers[eventType], handler)

	return nil
}

func (bus *MemoryEventBus) SubscribeAll(handler EventHandler) error {
	if bus.globalSubscribers == nil {
		bus.globalSubscribers = make([]EventHandler, 0)
	}

	bus.globalSubscribers = append(bus.globalSubscribers, handler)
	return nil
}
