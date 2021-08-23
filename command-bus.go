package goevents

import (
	"fmt"
	"log"
)

// A CommandBus is anything capable of command pub/sub.
type CommandBus interface {
	Handle(Command, CommandHandlerFn)
	Dispatch(Command) error
}

// The MemoryCommandBus is a CommandBus that dispatches commands to
// handlers registered locally in-memory.
type MemoryCommandBus struct {
	preHandlers  []CommandProcessorFn
	handlers     map[string]CommandHandlerFn
	postHandlers []CommandProcessorFn
}

var _ CommandBus = &MemoryCommandBus{}

// BeforeAny registers a handler that runs before any command is handled.
func (bus *MemoryCommandBus) BeforeAny(fn CommandProcessorFn) {
	bus.preHandlers = append(bus.preHandlers, fn)
}

// Handle registers the handler for a command type.
func (bus *MemoryCommandBus) Handle(commandType Command, handler CommandHandlerFn) {
	if bus.handlers == nil {
		bus.handlers = make(map[string]CommandHandlerFn)
	}

	bus.handlers[CommandName(commandType)] = handler
}

// AfterAny registers a handler that runs after any command is handled.
func (bus *MemoryCommandBus) AfterAny(fn CommandProcessorFn) {
	bus.postHandlers = append(bus.postHandlers, fn)
}

// Dispatch dispatches a command to the appropriate handler.
func (bus *MemoryCommandBus) Dispatch(cmd Command) error {
	data, err := cmd.Data()
	if err != nil {
		return err
	}

	log.Printf("[dispatch] %s %s\n", CommandName(cmd), string(data))

	for _, before := range bus.preHandlers {
		if err := before(cmd); err != nil {
			return err
		}
	}

	handler, ok := bus.handlers[CommandName(cmd)]
	if !ok {
		return fmt.Errorf("no registered command handlers for %q", CommandName(cmd))
	}

	if err := handler(cmd); err != nil {
		return err
	}

	for _, after := range bus.postHandlers {
		if err := after(cmd); err != nil {
			return err
		}
	}

	return nil
}
