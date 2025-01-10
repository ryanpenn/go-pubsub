package pubsub_test

import (
	"fmt"
	"pubsub"
)

const (
	TestEvent1 pubsub.EventType = iota + 1
	TestEvent2
	ValueChangedEvent
)

type ValueChanged struct {
	Old int
	New int
}

func ExampleEventManager() {
	em := pubsub.NewEventManager()
	em.OnEvent(0, func(event pubsub.Event) {
		fmt.Println(event.Type())
	})
	em.TriggerEvent(pubsub.EventType(0))

	em.OnEvent(TestEvent1, func(event pubsub.Event) {
		if _, ok := event.(pubsub.EventArger[string]); !ok {
			fmt.Println(event.Type(), "not a string event")
		}
	})
	em.TriggerEvent(TestEvent1)

	em.OnEvent(TestEvent2, func(event pubsub.Event) {
		if arg, ok := event.(pubsub.EventArger[string]); ok {
			fmt.Println(event.Type(), arg.Arg())
		}
	})
	em.TriggerEvent(pubsub.NewEventArg(TestEvent2, "hello world"))

	em.OnEvent(ValueChangedEvent, func(event pubsub.Event) {
		if arg, ok := event.(pubsub.EventArger[ValueChanged]); ok {
			fmt.Println(arg.Arg().Old)
			fmt.Println(arg.Arg().New)
		}
	})
	em.TriggerEvent(pubsub.NewEventArg(ValueChangedEvent, ValueChanged{10, 20}))

	// Output:
	// 0
	// 1 not a string event
	// 2 hello world
	// 10
	// 20
}
