package goevents

import (
	"fmt"
	"time"
)

// A CommandDispatcher is anything capable of command pub/sub.
type CommandDispatcher interface {
	On(Command, CommandHandler)
	Dispatch(Command) (*CommandResult, error)
}

// The LocalCommandDispatcher sis a synchronous, in-memory command dispatcher implementation.
type LocalCommandDispatcher struct {
	globalPreHandlers  []CommandProcessor
	handlers           map[string]CommandHandler
	globalPostHandlers []CommandProcessor
}

var _ CommandDispatcher = &LocalCommandDispatcher{}

// BeforeAny registers a handler that runs before any command is handled.
func (bus *LocalCommandDispatcher) BeforeAny(fn CommandProcessor) {
	bus.globalPreHandlers = append(bus.globalPreHandlers, fn)
}

// Handle registers the handler for a command type.
func (bus *LocalCommandDispatcher) On(commandType Command, handler CommandHandler) {
	if bus.handlers == nil {
		bus.handlers = make(map[string]CommandHandler)
	}

	bus.handlers[CommandName(commandType)] = handler
}

// AfterAny registers a handler that runs after any command is handled.
func (bus *LocalCommandDispatcher) AfterAny(fn CommandProcessor) {
	bus.globalPostHandlers = append(bus.globalPostHandlers, fn)
}

// Dispatch dispatches a command to the appropriate handler.
func (bus *LocalCommandDispatcher) Dispatch(cmd Command) (*CommandResult, error) {
	cmdName := CommandName(cmd)
	cmdCtx := &CommandContext{
		Type:         cmdName,
		DispatchedAt: time.Now(),
		Command:      cmd,
	}

	for _, before := range bus.globalPreHandlers {
		if err := before(cmdCtx); err != nil {
			return nil, err
		}
	}

	handler, ok := bus.handlers[cmdName]
	if !ok {
		return nil, fmt.Errorf("no registered command handlers for %q", cmdName)
	}

	result, err := handler(cmdCtx)
	if err != nil {
		return result, err
	}

	for _, after := range bus.globalPostHandlers {
		if err := after(cmdCtx); err != nil {
			return result, err
		}
	}

	return result, nil
}
