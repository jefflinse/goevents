package main

import (
	"fmt"

	"github.com/jefflinse/goevents"
)

func main() {
	// Create the event and command busses.
	events := goevents.MemoryEventBus{}
	commands := goevents.MemoryCommandBus{}

	// Create an event handler that will be called when *any* event is published.
	events.SubscribeAll(func(event *goevents.Event) error {
		fmt.Println("received event:", event.Type)
		return nil
	})

	// Create an event handler that will be called when the "SomethingHappened" event is published.
	events.Subscribe("SomethingHappened", func(event *goevents.Event) error {
		fmt.Println("something happened")
		return nil
	})

	// Create a command handler that will be called when the "DoSomething" command is dispatcheed.
	commands.Register("DoSomething", func(cmd *goevents.Command) error {
		fmt.Println("doing something")
		events.Publish(&goevents.Event{Type: "SomethingHappened"})
		return nil
	})

	commands.Dispatch(&goevents.Command{Type: "DoSomething"})
}
