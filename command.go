package goevents

import (
	"strings"
	"time"
)

type Command interface{}

func CommandName(c Command) string {
	return strings.TrimSuffix(typeName(c), "Command")
}

type CommandContext struct {
	Type         string
	DispatchedAt time.Time
	Command      Command
}

type CommandResult struct {
	Result interface{}
}

var EmptyCommandResult = &CommandResult{
	Result: nil,
}

// A CommandHandler is a function that handles a Command,
// returning a result and an error.
type CommandHandler func(cmd *CommandContext) (*CommandResult, error)

// A CommandContextProcessor is a function that pre- or pos-processes a Command,
// potentially modifying the CommandContext.
type CommandProcessor func(cmd *CommandContext) error
