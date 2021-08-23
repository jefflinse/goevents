package goevents

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
