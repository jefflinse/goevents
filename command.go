package goevents

import "time"

type Command interface{}

type CommandContext struct {
	Type         string
	DispatchedAt time.Time
	Command      Command
}

// A CommandHandlerFn is a function that handles a Command.
type CommandHandlerFn func(cmd *CommandContext) error
