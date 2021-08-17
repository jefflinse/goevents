package goevents

import (
	"fmt"
)

type Command struct {
	Type string
	UUID string
	Data []byte
}

type CommandHandler func(cmd *Command) error

type CommandBus interface {
	Register(commandType string, handler CommandHandler)
	Dispatch(cmd *Command) error
}

type MemoryCommandBus struct {
	handlers map[string]CommandHandler
}

func (bus *MemoryCommandBus) Register(commandType string, handler CommandHandler) {
	if bus.handlers == nil {
		bus.handlers = make(map[string]CommandHandler)
	}

	bus.handlers[commandType] = handler
}

func (bus *MemoryCommandBus) Dispatch(cmd *Command) error {
	handler, ok := bus.handlers[cmd.Type]
	if !ok {
		return fmt.Errorf("unhandled command: %q", cmd.Type)
	}

	return handler(cmd)
}
