package pubsub_test

import (
	"fmt"
	"pubsub"
)

func ExampleBroker() {
	b := pubsub.NewBroker[string]()
	us1 := b.Subscribe(func(topic string) {
		fmt.Println("sub1", topic)
	})
	defer us1.Unsubscribe()
	us2 := b.Subscribe(func(topic string) {
		fmt.Println("sub2", topic)
	})
	b.Publish("aaaaa")
	fmt.Println()
	us1.Unsubscribe()

	subject := pubsub.NewAutoUnsubscriber()
	b.Subscribe(func(topic string) {
		fmt.Println("sub3", topic)
	}).Bind(subject)
	b.Subscribe(func(topic string) {
		fmt.Println("sub4", topic)
	}).Bind(subject)
	b.Publish("bbbbb")
	fmt.Println()

	if f, ok := us2.(pubsub.UnsubscriberFunc); ok {
		f()
	}
	b.Publish("ccccc")

	subject.UnsubscribeAll()
	b.Publish("ddddd")

	// Output
	// sub1 aaaaa
	// sub2 aaaaa
	//
	// sub2 bbbbb
	// sub3 bbbbb
	// sub4 bbbbb
	//
	// sub3 ccccc
	// sub4 ccccc
}
