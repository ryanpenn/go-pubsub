package pubsub

import (
	"slices"
)

type Subscriber[T any] func(topic T)

type Unsubscriber interface {
	Unsubscribe()
	Bind(AutoUnsubscriber) // Bind binds the Unsubscriber to an AutoUnsubscriber.
}

type UnsubscriberFunc func()

func (f UnsubscriberFunc) Unsubscribe() {
	if f != nil {
		f()
	}
}

func (f UnsubscriberFunc) Bind(au AutoUnsubscriber) {
	if au != nil {
		au.Add(f)
	}
}

type AutoUnsubscriber interface {
	UnsubscribeAll()               // UnsubscribeAll unsubscribes all subscribers.
	Add(unsubscriber Unsubscriber) // Add adds an unsubscriber to the list.
}

type aus struct {
	unsubscribers []Unsubscriber
}

func (au *aus) UnsubscribeAll() {
	for _, u := range au.unsubscribers {
		if u != nil {
			u.Unsubscribe()
		}
	}
	au.unsubscribers = au.unsubscribers[:0]
}

func (au *aus) Add(unsubscriber Unsubscriber) {
	au.unsubscribers = append(au.unsubscribers, unsubscriber)
}

func NewAutoUnsubscriber() AutoUnsubscriber {
	return &aus{
		unsubscribers: make([]Unsubscriber, 0),
	}
}

type Broker[T any] struct {
	subscribers []*Subscriber[T]
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		subscribers: make([]*Subscriber[T], 0),
	}
}

func (b *Broker[T]) Subscribe(subscriber Subscriber[T]) Unsubscriber {
	b.subscribers = append(b.subscribers, &subscriber)
	return UnsubscriberFunc(func() {
		if b != nil {
			b.subscribers = slices.DeleteFunc(b.subscribers, func(s *Subscriber[T]) bool {
				return s == &subscriber
			})
		}
	})
}

func (b *Broker[T]) UnsubscribeAll() {
	b.subscribers = b.subscribers[:0]
}

func (b *Broker[T]) Publish(event T) {
	for _, handler := range b.subscribers {
		(*handler)(event)
	}
}
