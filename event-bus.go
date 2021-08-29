package goevents

import (
	"log"
	"time"
)

// An EventDispatcher is anything capable of event pub/sub.
type EventDispatcher interface {
	On(Event, EventHandler) error
	Publish(Event) error
}

// LocalEventDispatcher is an in-memory event dispatcher implementation.
type LocalEventDispatcher struct {
	globalPreHandlers  []EventHandler
	handlers           map[string][]EventHandler
	globalPostHandlers []EventHandler
}

var _ EventDispatcher = &LocalEventDispatcher{}

// Publish publishes an event, calling all registered handlers.
func (bus *LocalEventDispatcher) Publish(event Event) error {
	eventName := EventName(event)
	eventCtx := &EventContext{
		Type:         eventName,
		DispatchedAt: time.Now(),
		Event:        event,
	}

	log.Printf("[publish] %s %+v\n", eventName, eventCtx)

	for _, before := range bus.globalPreHandlers {
		if err := before(eventCtx); err != nil {
			return err
		}
	}

	for _, handle := range bus.handlers[EventName(event)] {
		if err := handle(eventCtx); err != nil {
			return err
		}
	}

	for _, after := range bus.globalPostHandlers {
		if err := after(eventCtx); err != nil {
			return err
		}
	}

	return nil
}

// BeforeAny registers a handler to be called before any event is dispatched.
// All pre-dispatch handlers are run in the order they are registered.
func (bus *LocalEventDispatcher) BeforeAny(handler EventHandler) error {
	if bus.globalPreHandlers == nil {
		bus.globalPreHandlers = make([]EventHandler, 0)
	}

	bus.globalPreHandlers = append(bus.globalPreHandlers, handler)
	return nil
}

// On registers a handler to be called when an event of the given type is published.
// All event handlers are run in the order they are registered.
func (bus *LocalEventDispatcher) On(eventType Event, handler EventHandler) error {
	if bus.handlers == nil {
		bus.handlers = make(map[string][]EventHandler)
	}

	eventName := EventName(eventType)

	if bus.handlers[eventName] == nil {
		bus.handlers[eventName] = make([]EventHandler, 0)
	}

	bus.handlers[eventName] = append(bus.handlers[eventName], handler)

	return nil
}

// AfterAny registers a handler to be called after any event is dispatched.
// All post-dispatch handlers are run in the order they are registered.
func (bus *LocalEventDispatcher) AfterAny(handler EventHandler) error {
	if bus.globalPostHandlers == nil {
		bus.globalPostHandlers = make([]EventHandler, 0)
	}

	bus.globalPostHandlers = append(bus.globalPostHandlers, handler)
	return nil
}
