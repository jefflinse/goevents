package main

import (
	"fmt"
	"time"

	"github.com/jefflinse/goevents"
)

type DoSomethingCommand struct {
	User string
}

type SomethingHappenedEvent struct {
	HappenedAt time.Time
	User       string
}

func main() {
	// Create the event and command busses.
	events := goevents.LocalEventDispatcher{}
	commands := goevents.LocalCommandDispatcher{}

	// Register an event handler that will be called before *any* event is handled.
	events.BeforeAny(func(e *goevents.EventContext) error {
		fmt.Printf("global event pre-handler (%s)\n  %+v\n", goevents.EventName(e.Event), e.Event)
		return nil
	})

	// Register an event handler that will be called when the "SomethingHappened" event is published.
	events.On(&SomethingHappenedEvent{}, func(e *goevents.EventContext) error {
		fmt.Printf("event handler (%s)\n  %+v\n", goevents.EventName(e.Event), e.Event)

		// Access event data by casting to known event type.
		// This is guaranteed to succeed because the event handler was registered using the event type name.
		event := e.Event.(*SomethingHappenedEvent)
		fmt.Printf("  user: %s\n", event.User)
		return nil
	})

	// Register an event handler that wsill be called after *any* event is handled.
	events.AfterAny(func(e *goevents.EventContext) error {
		fmt.Printf("global event post-handler (%s)\n  %+v\n", goevents.EventName(e.Event), e.Event)
		return nil
	})

	commands.BeforeAny(func(c *goevents.CommandContext) error {
		fmt.Printf("global command pre-handler (%s)\n  %+v\n", goevents.CommandName(c.Command), c.Command)
		return nil
	})

	// Create a command handler that will be called when the "DoSomething" command is dispatcheed.
	commands.On(&DoSomethingCommand{}, func(c *goevents.CommandContext) (*goevents.CommandResult, error) {
		fmt.Printf("command handler (%s)\n  %+v\n", goevents.CommandName(c.Command), c.Command)

		// Access command data by casting to known command type.
		// This is guaranteed to succeed because the command handler was registered using the command type name.
		cmd := c.Command.(*DoSomethingCommand)
		fmt.Printf("  user: %s\n", cmd.User)
		if err := events.Dispatch(&SomethingHappenedEvent{
			User:       cmd.User,
			HappenedAt: time.Now(),
		}); err != nil {
			// If the event could not be dispatched, handle it here,
			// such as reverting the changes caused by the command.
			return nil, err
		}

		return goevents.EmptyCommandResult, nil
	})

	commands.AfterAny(func(c *goevents.CommandContext) error {
		fmt.Printf("global command post-handler (%s)\n  %+v\n", goevents.CommandName(c.Command), c.Command)
		return nil
	})

	result, err := commands.Dispatch(&DoSomethingCommand{User: "Jeff"})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("result: %+v\n", result)
	}
}
