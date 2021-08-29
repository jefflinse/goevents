package goevents

import (
	"time"
)

// An EventDispatcher is anything capable of event pub/sub.
type EventDispatcher interface {
	On(Event, EventHandler) error
	Dispatch(Event) error
}

// LocalEventDispatcher is a synchronous, in-memory event dispatcher implementation.
type LocalEventDispatcher struct {
	globalPreHandlers  []EventHandler
	handlers           map[string][]EventHandler
	globalPostHandlers []EventHandler
}

var _ EventDispatcher = &LocalEventDispatcher{}

// BeforeAny registers a handler to be called before any event is dispatched.
// All pre-dispatch handlers are run in the order they are registered.
func (bus *LocalEventDispatcher) BeforeAny(handler EventHandler) error {
	bus.globalPreHandlers = append(bus.globalPreHandlers, handler)
	return nil
}

// On registers a handler to be called when an event of the given type is dispatched.
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
	bus.globalPostHandlers = append(bus.globalPostHandlers, handler)
	return nil
}

// Dispatch dispatches an event, calling all registered handlers.
func (bus *LocalEventDispatcher) Dispatch(event Event) error {
	eventName := EventName(event)
	eventCtx := &EventContext{
		Type:         eventName,
		DispatchedAt: time.Now(),
		Event:        event,
	}

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
