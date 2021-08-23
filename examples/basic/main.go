package main

import (
	"encoding/json"
	"fmt"

	"github.com/jefflinse/goevents"
)

type DoSomethingCommand struct {
	goevents.JSONCommand
	User   string
	Before string
	After  string
}

func (c *DoSomethingCommand) Name() string {
	return "DoSomething"
}

func (c *DoSomethingCommand) Data() ([]byte, error) {
	return json.Marshal(c)
}

type SomethingHappenedEvent struct {
	goevents.JSONEvent
	Before string
	After  string
}

func main() {
	// Create the event and command busses.
	events := goevents.MemoryEventBus{}
	commands := goevents.DefaultCommandBus{}

	// Register an event handler that will be called before *any* event is handled.
	events.BeforeAny(func(e goevents.Event) error {
		data, err := e.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside global pre-event handler! (type: %s, data: %s)\n", goevents.EventName(e), string(data))
		return nil
	})

	// Register an event handler that will be called when the "SomethingHappened" event is published.
	events.On("SomethingHappened", func(e goevents.Event) error {
		event := e.(*SomethingHappenedEvent)
		data, err := event.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside %s event handler! (data: %s)\n", goevents.EventName(e), string(data))

		fmt.Printf("  before: %s\n", event.Before)
		fmt.Printf("  after:  %s\n", event.After)
		return nil
	})

	// Register an event handler that wsill be called after *any* event is handled.
	events.AfterAny(func(e goevents.Event) error {
		data, err := e.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside global post-event handler! (type: %s, data: %+v)\n", goevents.EventName(e), string(data))
		return nil
	})

	commands.BeforeAny(func(c goevents.Command) error {
		data, err := c.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside global pre-command handler! (type: %s, data: %s)\n", goevents.CommandName(c), string(data))
		return nil
	})

	// Create a command handler that will be called when the "DoSomething" command is dispatcheed.
	commands.Handle("DoSomething", func(c goevents.Command) (goevents.CommandResult, error) {
		cmd := c.(*DoSomethingCommand)
		fmt.Printf("inside DoSomething command handler!\n")
		events.Publish(&SomethingHappenedEvent{Before: cmd.Before, After: cmd.After})
		return nil, nil
	})

	commands.AfterAny(func(c goevents.Command) error {
		data, err := c.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside global post-command handler! (type: %s, data: %s)\n", goevents.CommandName(c), string(data))
		return nil
	})

	commands.Dispatch(&DoSomethingCommand{User: "Jeff", Before: "before", After: "after"})
}
