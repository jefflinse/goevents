package goevents

import (
	"log"
	"time"
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

// Publish publishes an event, calling all registered handlers.
func (bus *MemoryEventBus) Publish(event Event) error {
	eventName := EventName(event)
	eventCtx := &EventContext{
		Type:         eventName,
		DispatchedAt: time.Now(),
		Event:        event,
	}

	log.Printf("[publish] %s %+v\n", eventName, eventCtx)

	for _, before := range bus.globalPreSubscribers {
		if err := before(eventCtx); err != nil {
			return err
		}
	}

	for _, handle := range bus.subscribers[EventName(event)] {
		if err := handle(eventCtx); err != nil {
			return err
		}
	}

	for _, after := range bus.globalPostSubscribers {
		if err := after(eventCtx); err != nil {
			return err
		}
	}

	return nil
}

// BeforeAny registers a handler to be called before any event is published.
func (bus *MemoryEventBus) BeforeAny(handler EventHandler) error {
	if bus.globalPreSubscribers == nil {
		bus.globalPreSubscribers = make([]EventHandler, 0)
	}

	bus.globalPreSubscribers = append(bus.globalPreSubscribers, handler)
	return nil
}

// On registers a handler to be called when an event of the given type is published.
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

// AfterAny registers a handler to be called after any event is published.
func (bus *MemoryEventBus) AfterAny(handler EventHandler) error {
	if bus.globalPostSubscribers == nil {
		bus.globalPostSubscribers = make([]EventHandler, 0)
	}

	bus.globalPostSubscribers = append(bus.globalPostSubscribers, handler)
	return nil
}
