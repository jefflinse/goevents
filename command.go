package goevents

type Command interface {
	Data() ([]byte, error)
}

// A CommandHandlerFn is a function that handles a Command.
type CommandHandlerFn func(cmd Command) error

// A CommandProcessorFn is a function that performs an action before or after handling a Command.
type CommandProcessorFn func(cmd Command) error
