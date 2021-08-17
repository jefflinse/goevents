package main

import (
	"fmt"

	"github.com/jefflinse/goevents"
)

func main() {
	// Create the event and command busses.
	events := goevents.MemoryEventBus{}
	commands := goevents.DefaultCommandDispatcher{}

	// Register an event handler that will be called when *any* event is published.
	events.OnAll(func(event *goevents.Event) error {
		fmt.Println("received event:", event.Type)
		return nil
	})

	// Register an event handler that will be called when the "SomethingHappened" event is published.
	events.On("SomethingHappened", func(event *goevents.Event) error {
		fmt.Println("something happened")
		return nil
	})

	// Create a command handler that will be called when the "DoSomething" command is dispatcheed.
	commands.Handle("DoSomething", func(cmd *goevents.Command) error {
		events.Publish(&goevents.Event{Type: "SomethingHappened"})
		return nil
	})

	commands.Dispatch(&goevents.Command{Type: "DoSomething"})
}
