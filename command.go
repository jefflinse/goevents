package goevents

type Command interface {
	Name() string
	Data() ([]byte, error)
}

type CommandResult interface {
	Data() []byte
}

// A CommandHandlerFn is a function that handles a Command.
type CommandHandlerFn func(cmd Command) (CommandResult, error)

// A CommandProcessorFn is a function that performs an action before or after handling a Command.
type CommandProcessorFn func(cmd Command) error

type DefaultCommand struct {
	CommandName string
	CommandData []byte
}

func (c *DefaultCommand) Name() string {
	return c.CommandName
}

func (c *DefaultCommand) Data() ([]byte, error) {
	return c.CommandData, nil
}
