package goevents

import (
	"fmt"
	"log"
)

type Command interface {
	Type() string
	Data() []byte
}

type CommandHandler func(cmd Command) error

type CommandDispatcher interface {
	Handle(commandType string, handler CommandHandler)
	Dispatch(cmd Command) error
}

type DefaultCommandDispatcher struct {
	handlers map[string]CommandHandler
}

var _ CommandDispatcher = &DefaultCommandDispatcher{}

func (dispatcher *DefaultCommandDispatcher) Handle(commandType string, handler CommandHandler) {
	if dispatcher.handlers == nil {
		dispatcher.handlers = make(map[string]CommandHandler)
	}

	dispatcher.handlers[commandType] = handler
}

func (dispatcher *DefaultCommandDispatcher) Dispatch(cmd Command) error {
	log.Printf("[command] %s {%s}", cmd.Type(), string(cmd.Data()))

	handler, ok := dispatcher.handlers[cmd.Type()]
	if !ok {
		return fmt.Errorf("unhandled command: %q", cmd.Type())
	}

	return handler(cmd)
}
