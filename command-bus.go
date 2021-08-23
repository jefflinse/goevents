package goevents

import (
	"fmt"
	"log"
)

// A CommandBus is anything capable of command pub/sub.
type CommandBus interface {
	Handle(string, CommandHandlerFn)
	Dispatch(Command) (CommandResult, error)
}

// The DefaultCommandBus is a CommandBus that dispatches commands to
// handlers registered locally in-memory.
type DefaultCommandBus struct {
	preHandlers  []CommandProcessorFn
	handlers     map[string]CommandHandlerFn
	postHandlers []CommandProcessorFn
}

var _ CommandBus = &DefaultCommandBus{}

// BeforeAny registers a handler that runs before any command is handled.
func (bus *DefaultCommandBus) BeforeAny(fn CommandProcessorFn) {
	bus.preHandlers = append(bus.preHandlers, fn)
}

// Handle registers the handler for a command type.
func (bus *DefaultCommandBus) Handle(cmdName string, handler CommandHandlerFn) {
	if bus.handlers == nil {
		bus.handlers = make(map[string]CommandHandlerFn)
	}

	bus.handlers[cmdName] = handler
}

// AfterAny registers a handler that runs after any command is handled.
func (bus *DefaultCommandBus) AfterAny(fn CommandProcessorFn) {
	bus.postHandlers = append(bus.postHandlers, fn)
}

// Dispatch dispatches a command to the appropriate handler.
func (bus *DefaultCommandBus) Dispatch(cmd Command) (CommandResult, error) {
	data, err := cmd.Data()
	if err != nil {
		return nil, err
	}

	log.Printf("[dispatch] %s %s\n", cmd.UID(), string(data))

	for _, before := range bus.preHandlers {
		if err := before(cmd); err != nil {
			return nil, err
		}
	}

	handler, ok := bus.handlers[cmd.UID()]
	if !ok {
		return nil, fmt.Errorf("no registered command handlers for %q", cmd.UID())
	}

	result, err := handler(cmd)
	if err != nil {
		return nil, err
	}

	for _, after := range bus.postHandlers {
		if err := after(cmd); err != nil {
			return nil, err
		}
	}

	return result, nil
}
