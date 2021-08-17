package goevents

import "log"

type Event struct {
	Type string
	UUID string
	Data []byte
}

type EventHandler func(event *Event) error

type EventBus interface {
	Publish(event *Event) error
	Subscribe(eventType string, handler EventHandler) error
}

type MemoryEventBus struct {
	globalSubscribers []EventHandler
	subscribers       map[string][]EventHandler
}

func (bus *MemoryEventBus) Publish(event *Event) error {
	log.Printf("[event] %s {%s}\n", event.Type, string(event.Data))

	for _, globalHandler := range bus.globalSubscribers {
		if err := globalHandler(event); err != nil {
			return err
		}
	}

	for _, handler := range bus.subscribers[event.Type] {
		if err := handler(event); err != nil {
			return err
		}
	}

	return nil
}

func (bus *MemoryEventBus) On(eventType string, handler EventHandler) error {
	if bus.subscribers == nil {
		bus.subscribers = make(map[string][]EventHandler)
	}

	if bus.subscribers[eventType] == nil {
		bus.subscribers[eventType] = make([]EventHandler, 0)
	}

	bus.subscribers[eventType] = append(bus.subscribers[eventType], handler)

	return nil
}

func (bus *MemoryEventBus) OnAll(handler EventHandler) error {
	if bus.globalSubscribers == nil {
		bus.globalSubscribers = make([]EventHandler, 0)
	}

	bus.globalSubscribers = append(bus.globalSubscribers, handler)
	return nil
}
