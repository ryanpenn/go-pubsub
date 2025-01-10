# go-pubsub

A simple pubsub library for Go.

## pub-sub model

The pubsub model is a simple messaging pattern where publishers (publishers) send messages to a topic (topic) and subscribers (subscribers) receive messages from the topic. In this library, we use the `Broker` interface to represent the publisher and the `Subscriber[T any]` type to represent the subscriber.

- Subscribe to a topic
  ```go
    b := pubsub.NewBroker[string]()
    b.Subscribe(func(topic string) {
        fmt.Println("received", topic)
    })
    b.Publish("hello")
    // Output:
    // received hello
  ```
- Unsubscribe from a topic
  ```go
    b := pubsub.NewBroker[string]()
    unsub := b.Subscribe(func(topic string) {
        fmt.Println("received", topic)
    })
    b.Publish("hello")
    unsub.Unsubscribe()
    b.Publish("world")
    // Output:
    // received hello
  ```
- Bind an Subject's lifecyle
  ```go
    autoUnsub := pubsub.NewAutoUnsubscribe()
    b := pubsub.NewBroker[string]()
    b.Subscribe(func(topic string) {
        fmt.Println("received", topic)
    }).Bind(autoUnsub)
    b.Subscribe(func(topic string) {
        fmt.Println("message", topic)
    }).Bind(autoUnsub)
    b.Publish("hello")
    // Output:
    // received hello
    // message hello
    autoUnsub.UnsubscribeAll()
    b.Publish("world")
    // Output:
  ```
- Pub-sub example
  ```go
    type Option func(*obj)

    func WithBroker(broker *pubsub.Broker[string]) Option {
      return func(o *obj) {
        broker.Subscribe(o.OnMessage).Bind(o)
      }
    }

    func WithName(name string) Option {
      return func(o *obj) {
        o.name = name
      }
    }

    type obj struct {
      pubsub.AutoUnsubscriber
      name string
    }

    func NewObj(opt ...Option) *obj {
      o := &obj{
        AutoUnsubscriber: pubsub.NewAutoUnsubscriber(),
      }
      for _, f := range opt {
        f(o)
      }
      return o
    }

    func (o *obj) OnMessage(msg string) {
      fmt.Println(o.name, "received message:", msg)
    }

    func (o *obj) Release() {
      o.UnsubscribeAll()
    }

    func main() {
      var broker = pubsub.NewBroker[string]()

      obj1 := NewObj(WithName("obj1"), WithBroker(broker))
      obj2 := NewObj(WithName("obj2"), WithBroker(broker))

      broker.Publish("Hello")
      obj1.Release()

      broker.Publish("World")
      obj2.Release()

      broker.Publish("!")

      // Output:
      // obj1 received message: Hello
      // obj2 received message: Hello
      // obj2 received message: World
    }
  ```

## event

The `event` provides a simple event system that allows you to subscribe to events and trigger them.

- Subscribe to an event
  ```go
    em := pubsub.NewEventManager()
	em.OnEvent(0, func(event pubsub.Event) {
		fmt.Println(event.Type())
	})
	em.TriggerEvent(pubsub.EventType(0))
    // Output:
	// 0
  ```
- Unsubscribe from an event
  ```go
    em := pubsub.NewEventManager()
	unsub := em.OnEvent(0, func(event pubsub.Event) {
		fmt.Println(event.Type())
	})
	em.TriggerEvent(pubsub.EventType(0))
	unsub.Unsubscribe()
	em.TriggerEvent(pubsub.EventType(0))
    // Output:
	// 0
  ```
- Event with data
  ```go
    eType := pubsub.EventType(1)
    em := pubsub.NewEventManager()
	em.OnEvent(eType, func(event pubsub.Event) {
        fmt.Println("type check", event.Type()==eType)
		if arg, ok := event.(pubsub.EventArger[string]); ok {
			fmt.Println("arg", event.Type(), arg.Arg())
		}
	})
	em.TriggerEvent(pubsub.NewEventArg(eType, "hello world"))
    // Output:
	// type check true
    // arg 1 hello world
  ```
- Custom event arg
  ```go
    type ValueChanged struct {
        Old int
        New int
    }
    em := pubsub.NewEventManager()
	em.OnEvent(pubsub.EventType(2), func(event pubsub.Event) {
		if arg, ok := event.(pubsub.EventArger[ValueChanged]); ok {
			fmt.Println("value changed", arg.Arg().Old, "->", arg.Arg().New)
		}
	})
	em.TriggerEvent(pubsub.NewEventArg(pubsub.EventType(2), ValueChanged{1, 2}))
    // Output:
	// value changed 1 -> 2
  ```
