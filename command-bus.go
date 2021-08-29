package goevents

import (
	"fmt"
	"time"
)

// A CommandDispatcher is anything capable of command pub/sub.
type CommandDispatcher interface {
	On(Command, CommandHandlerFn)
	Dispatch(Command) error
}

// The LocalCommandDispatcher is a synchronous, in-memory command dispatcher implementation.
type LocalCommandDispatcher struct {
	globalPreHandlers  []CommandHandlerFn
	handlers           map[string]CommandHandlerFn
	globalPostHandlers []CommandHandlerFn
}

var _ CommandDispatcher = &LocalCommandDispatcher{}

// BeforeAny registers a handler that runs before any command is handled.
func (bus *LocalCommandDispatcher) BeforeAny(fn CommandHandlerFn) {
	bus.globalPreHandlers = append(bus.globalPreHandlers, fn)
}

// Handle registers the handler for a command type.
func (bus *LocalCommandDispatcher) On(commandType Command, handler CommandHandlerFn) {
	if bus.handlers == nil {
		bus.handlers = make(map[string]CommandHandlerFn)
	}

	bus.handlers[CommandName(commandType)] = handler
}

// AfterAny registers a handler that runs after any command is handled.
func (bus *LocalCommandDispatcher) AfterAny(fn CommandHandlerFn) {
	bus.globalPostHandlers = append(bus.globalPostHandlers, fn)
}

// Dispatch dispatches a command to the appropriate handler.
func (bus *LocalCommandDispatcher) Dispatch(cmd Command) error {
	cmdName := CommandName(cmd)
	cmdCtx := &CommandContext{
		Type:         cmdName,
		DispatchedAt: time.Now(),
		Command:      cmd,
	}

	for _, before := range bus.globalPreHandlers {
		if err := before(cmdCtx); err != nil {
			return err
		}
	}

	handler, ok := bus.handlers[cmdName]
	if !ok {
		return fmt.Errorf("no registered command handlers for %q", cmdName)
	}

	if err := handler(cmdCtx); err != nil {
		return err
	}

	for _, after := range bus.globalPostHandlers {
		if err := after(cmdCtx); err != nil {
			return err
		}
	}

	return nil
}
