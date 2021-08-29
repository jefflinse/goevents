package main

import (
	"encoding/json"
	"fmt"

	"github.com/jefflinse/goevents"
)

type DoSomethingCommand struct {
	User string
}

func (c *DoSomethingCommand) Data() ([]byte, error) {
	return json.Marshal(c)
}

type SomethingHappenedEvent struct {
	User string
}

func (e *SomethingHappenedEvent) Data() ([]byte, error) {
	return json.Marshal(e)
}

func main() {
	// Create the event and command busses.
	events := goevents.LocalEventDispatcher{}
	commands := goevents.LocalCommandDispatcher{}

	// Register an event handler that will be called before *any* event is handled.
	events.BeforeAny(func(e *goevents.EventContext) error {
		fmt.Printf("global event pre-handler (%s)\n  %+v\n", goevents.EventName(e.Event), e)
		return nil
	})

	// Register an event handler that will be called when the "SomethingHappened" event is published.
	events.On(&SomethingHappenedEvent{}, func(e *goevents.EventContext) error {
		fmt.Printf("event handler (%s)\n  %+v\n", goevents.EventName(e.Event), e)

		// access event data by casting to known event type
		event := e.Event.(*SomethingHappenedEvent)
		fmt.Printf("  user: %s\n", event.User)
		return nil
	})

	// Register an event handler that wsill be called after *any* event is handled.
	events.AfterAny(func(e *goevents.EventContext) error {
		fmt.Printf("global event post-handler (%s)\n  %+v\n", goevents.EventName(e.Event), e)
		return nil
	})

	commands.BeforeAny(func(c *goevents.CommandContext) error {
		fmt.Printf("global command pre-handler (%s)\n  %+v\n", goevents.CommandName(c.Command), c)
		return nil
	})

	// Create a command handler that will be called when the "DoSomething" command is dispatcheed.
	commands.On(&DoSomethingCommand{}, func(c *goevents.CommandContext) error {
		fmt.Printf("command handler (%s)\n  %+v\n", goevents.CommandName(c.Command), c)

		cmd := c.Command.(*DoSomethingCommand)
		if err := events.Dispatch(&SomethingHappenedEvent{User: cmd.User}); err != nil {
			// If the event could not be published, handle it here,
			// such as reverting the changes caused by the command.
			return err
		}

		return nil
	})

	commands.AfterAny(func(c *goevents.CommandContext) error {
		fmt.Printf("global command post-handler (%s)\n  %+v\n", goevents.CommandName(c.Command), c)
		return nil
	})

	commands.Dispatch(&DoSomethingCommand{User: "Jeff"})
}
