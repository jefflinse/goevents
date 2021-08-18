package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jefflinse/goevents"
)

type DoSomethingCommand struct {
	Before string
	After  string
}

func (c *DoSomethingCommand) Type() string {
	return "DoSomething"
}

func (c *DoSomethingCommand) Data() []byte {
	b, _ := json.Marshal(c)
	return b
}

type SomethingHappenedEvent struct {
	Before string
	After  string
	When   time.Time
}

func (c *SomethingHappenedEvent) Type() string {
	return "SomethingHappened"
}

func (c *SomethingHappenedEvent) Data() []byte {
	b, _ := json.Marshal(c)
	return b
}

func (c *SomethingHappenedEvent) OccurredAt() time.Time {
	return c.When
}

func main() {
	// Create the event and command busses.
	events := goevents.MemoryEventBus{}
	commands := goevents.DefaultCommandBus{}

	// Register an event handler that will be called when *any* event is published.
	events.SubscribeAll(func(e goevents.Event) error {
		fmt.Printf("inside global event handler! (type: %s, data: %s)\n", e.Type(), string(e.Data()))
		return nil
	})

	// Register an event handler that will be called when the "SomethingHappened" event is published.
	events.Subscribe("SomethingHappened", func(e goevents.Event) error {
		fmt.Printf("inside %s event handler! (data: %s)\n", e.Type(), string(e.Data()))
		return nil
	})

	// Create a command handler that will be called when the "DoSomething" command is dispatcheed.
	commands.Handle("DoSomething", func(c goevents.Command) error {
		cmd := c.(*DoSomethingCommand)
		fmt.Printf("inside %s command handler!\n", cmd.Type())
		events.Publish(&SomethingHappenedEvent{Before: cmd.Before, After: cmd.After, When: time.Now()})
		return nil
	})

	commands.Dispatch(&DoSomethingCommand{Before: "before", After: "after"})
}
