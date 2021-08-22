package goevents

import (
	"log"
	"time"
)

type Event interface {
	At() time.Time
	Data() []byte
	Type() string
}

type EventHandler func(event Event) error

type BasicEvent struct {
	ETime time.Time
	EData []byte
}

func (e *BasicEvent) At() time.Time {
	return e.ETime
}

func (e *BasicEvent) Data() []byte {
	return e.EData
}

func (e *BasicEvent) Type() string {
	return getEventType(e)
}

var _ Event = &BasicEvent{}

// An EventBus is anything capable of event pub/sub.
type EventBus interface {
	Publish(event Event) error
	Subscribe(eventType Event, handler EventHandler) error
	SubscribeAll(handler EventHandler) error
}

// MemoryEventBus is an in-memory event bus implementation.
type MemoryEventBus struct {
	globalSubscribers []EventHandler
	subscribers       map[string][]EventHandler
}

var _ EventBus = &MemoryEventBus{}

func (bus *MemoryEventBus) Publish(event Event) error {
	eventType := getEventType(event)
	log.Printf("[event] %s {%s}\n", eventType, string(event.Data()))

	for _, globalHandler := range bus.globalSubscribers {
		if err := globalHandler(event); err != nil {
			return err
		}
	}

	for _, handler := range bus.subscribers[eventType] {
		if err := handler(event); err != nil {
			return err
		}
	}

	return nil
}

func (bus *MemoryEventBus) Subscribe(eventType Event, handler EventHandler) error {
	eventTypeName := getEventType(eventType)

	if bus.subscribers == nil {
		bus.subscribers = make(map[string][]EventHandler)
	}

	if bus.subscribers[eventTypeName] == nil {
		bus.subscribers[eventTypeName] = make([]EventHandler, 0)
	}

	bus.subscribers[eventTypeName] = append(bus.subscribers[eventTypeName], handler)

	return nil
}

func (bus *MemoryEventBus) SubscribeAll(handler EventHandler) error {
	if bus.globalSubscribers == nil {
		bus.globalSubscribers = make([]EventHandler, 0)
	}

	bus.globalSubscribers = append(bus.globalSubscribers, handler)
	return nil
}
