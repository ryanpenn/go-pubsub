package pubsub

// EventType is a type of event, which implements the Event interface
type EventType int

func (t EventType) Type() EventType {
	return t
}

// Event is an interface for all events
type Event interface {
	Type() EventType
}

// EventArger is an interface for events with arguments
type EventArger[T any] interface {
	Event
	Arg() T
}

// eventArg is an EventArger generic implementation
type eventArg[T any] struct {
	eType EventType
	data  T
}

func (e eventArg[T]) Type() EventType {
	return e.eType
}

func (e eventArg[T]) Arg() T {
	return e.data
}

func NewEventArg[T any](t EventType, data T) EventArger[T] {
	return eventArg[T]{
		eType: t,
		data:  data,
	}
}

// EventManager is a manager for events
type EventManager struct {
	subs map[EventType]*Broker[Event]
}

func NewEventManager() *EventManager {
	return &EventManager{
		subs: make(map[EventType]*Broker[Event]),
	}
}

func (g *EventManager) OnEvent(event EventType, fn Subscriber[Event]) Unsubscriber {
	if _, ok := g.subs[event]; !ok {
		g.subs[event] = NewBroker[Event]()
	}
	return g.subs[event].Subscribe(fn)
}

func (g *EventManager) TriggerEvent(event Event) {
	if _, ok := g.subs[event.Type()]; ok {
		g.subs[event.Type()].Publish(event)
	}
}
