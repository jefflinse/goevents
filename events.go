package goevents

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
	globalSubscribes []EventHandler
	subscribers      map[string][]EventHandler
}

func (bus *MemoryEventBus) Publish(event *Event) error {
	for _, handler := range bus.globalSubscribes {
		if err := handler(event); err != nil {
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
	if bus.globalSubscribes == nil {
		bus.globalSubscribes = make([]EventHandler, 0)
	}

	bus.globalSubscribes = append(bus.globalSubscribes, handler)
	return nil
}
