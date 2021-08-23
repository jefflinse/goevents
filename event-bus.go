package goevents

import (
	"log"
)

// An EventPublisher is anything that can publish events.
type EventPublisher interface {
	Publish(Event) error
}

// An EventSubscriber is anything that can subscribe to events.
type EventSubscriber interface {
	On(Event, EventHandler) error
}

// An EventBus is anything capable of event pub/sub.
type EventBus interface {
	EventPublisher
	EventSubscriber
}

// MemoryEventBus is an in-memory event bus implementation.
type MemoryEventBus struct {
	globalPreSubscribers  []EventHandler
	subscribers           map[string][]EventHandler
	globalPostSubscribers []EventHandler
}

var _ EventBus = &MemoryEventBus{}

func (bus *MemoryEventBus) Publish(event Event) error {
	data, err := event.Data()
	if err != nil {
		return err
	}

	eventName := EventName(event)
	log.Printf("[publish] %s %s\n", eventName, string(data))

	for _, before := range bus.globalPreSubscribers {
		if err := before(event); err != nil {
			return err
		}
	}

	for _, handle := range bus.subscribers[EventName(event)] {
		if err := handle(event); err != nil {
			return err
		}
	}

	for _, after := range bus.globalPostSubscribers {
		if err := after(event); err != nil {
			return err
		}
	}

	return nil
}

func (bus *MemoryEventBus) BeforeAny(handler EventHandler) error {
	if bus.globalPreSubscribers == nil {
		bus.globalPreSubscribers = make([]EventHandler, 0)
	}

	bus.globalPreSubscribers = append(bus.globalPreSubscribers, handler)
	return nil
}

func (bus *MemoryEventBus) On(eventType Event, handler EventHandler) error {
	if bus.subscribers == nil {
		bus.subscribers = make(map[string][]EventHandler)
	}

	eventName := EventName(eventType)

	if bus.subscribers[eventName] == nil {
		bus.subscribers[eventName] = make([]EventHandler, 0)
	}

	bus.subscribers[eventName] = append(bus.subscribers[eventName], handler)

	return nil
}

func (bus *MemoryEventBus) AfterAny(handler EventHandler) error {
	if bus.globalPostSubscribers == nil {
		bus.globalPostSubscribers = make([]EventHandler, 0)
	}

	bus.globalPostSubscribers = append(bus.globalPostSubscribers, handler)
	return nil
}
