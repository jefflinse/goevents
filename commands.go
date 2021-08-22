package goevents

import (
	"fmt"
	"log"
)

type Command interface {
	Type() string
	Data() []byte
}

// A CommandHandler is a function that handles a Command.
type CommandHandler func(cmd Command) error

// A CommandBus is anything capable of command pub/sub.
type CommandBus interface {
	Handle(commandType string, handler CommandHandler)
	Dispatch(cmd Command) error
}

// The DefaultCommandBus is a CommandBus that dispatches commands to
// handlesr registered locally in-memory.
type DefaultCommandBus struct {
	handlers map[string]CommandHandler
}

var _ CommandBus = &DefaultCommandBus{}

// Handle registers a handler for a command type.
func (dispatcher *DefaultCommandBus) Handle(commandType string, handler CommandHandler) {
	if dispatcher.handlers == nil {
		dispatcher.handlers = make(map[string]CommandHandler)
	}

	dispatcher.handlers[commandType] = handler
}

// Dispatch dispatches a command to the appropriate handler.
func (dispatcher *DefaultCommandBus) Dispatch(cmd Command) error {
	log.Printf("[command] %s {%s}", cmd.Type(), string(cmd.Data()))

	handler, ok := dispatcher.handlers[cmd.Type()]
	if !ok {
		return fmt.Errorf("unhandled command: %q", cmd.Type())
	}

	return handler(cmd)
}
