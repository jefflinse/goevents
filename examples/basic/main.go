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
	events := goevents.MemoryEventBus{}
	commands := goevents.MemoryCommandBus{}

	// Register an event handler that will be called before *any* event is handled.
	events.BeforeAny(func(e goevents.Event) error {
		data, err := e.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside global pre-event handler!\n  (%s) %s\n", goevents.EventName(e), string(data))
		return nil
	})

	// Register an event handler that will be called when the "SomethingHappened" event is published.
	events.On(&SomethingHappenedEvent{}, func(e goevents.Event) error {
		event := e.(*SomethingHappenedEvent)
		data, err := e.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside %s event handler!\n  %s\n", goevents.EventName(e), string(data))

		fmt.Printf("  user: %s\n", event.User)
		return nil
	})

	// Register an event handler that wsill be called after *any* event is handled.
	events.AfterAny(func(e goevents.Event) error {
		data, err := e.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside global post-event handler!\n  (%s) %s\n", goevents.EventName(e), string(data))
		return nil
	})

	commands.BeforeAny(func(c goevents.Command) error {
		data, err := c.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside global pre-command handler\n  (%s) %s\n", goevents.CommandName(c), string(data))
		return nil
	})

	// Create a command handler that will be called when the "DoSomething" command is dispatcheed.
	commands.Handle(&DoSomethingCommand{}, func(c goevents.Command) error {
		data, err := c.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside %s command handler!\n  %s\n", goevents.CommandName(c), string(data))

		cmd := c.(*DoSomethingCommand)
		if err := events.Publish(&SomethingHappenedEvent{User: cmd.User}); err != nil {
			// If the event could not be published, handle it here,
			// such as reverting the changes caused by the command.
			return err
		}

		return nil
	})

	commands.AfterAny(func(c goevents.Command) error {
		data, err := c.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside global post-command handler\n  (%s) %s\n", goevents.CommandName(c), string(data))
		return nil
	})

	commands.Dispatch(&DoSomethingCommand{User: "Jeff"})
}
