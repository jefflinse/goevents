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

// A CommandHandler is a function that handles a Command.
type CommandHandler func(cmd *CommandContext) error
