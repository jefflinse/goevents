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

		fmt.Printf("inside global pre-event handler! (type: %s, data: %s)\n", goevents.EventName(e), string(data))
		return nil
	})

	// Register an event handler that will be called when the "SomethingHappened" event is published.
	events.On(&SomethingHappenedEvent{}, func(e goevents.Event) error {
		event := e.(*SomethingHappenedEvent)
		data, err := e.Data()
		if err != nil {
			return err
		}

		fmt.Printf("inside %s event handler! (data: %s)\n", goevents.EventName(e), string(data))

		fmt.Printf("  user: %s\n", event.User)
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
	commands.Handle(&DoSomethingCommand{}, func(c goevents.Command) (goevents.CommandResult, error) {
		cmd := c.(*DoSomethingCommand)
		cmdName := goevents.CommandName(c)
		fmt.Printf("inside %s command handler!\n", cmdName)

		defer func() {
			events.Publish(&SomethingHappenedEvent{User: cmd.User})
		}()

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

	commands.Dispatch(&DoSomethingCommand{User: "Jeff"})
}
