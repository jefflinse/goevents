package goevents

import "encoding/json"

type Command interface {
	Data() ([]byte, error)
}

type CommandResult interface {
	Data() []byte
}

// A CommandHandlerFn is a function that handles a Command.
type CommandHandlerFn func(cmd Command) (CommandResult, error)

// A CommandProcessorFn is a function that performs an action before or after handling a Command.
type CommandProcessorFn func(cmd Command) error

// A JSONCommand is a command whose data is simply the command object marshalled as JSON.
type JSONCommand struct{}

func (c *JSONCommand) Data() ([]byte, error) {
	return json.Marshal(c)
}
