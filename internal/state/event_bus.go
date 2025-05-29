package state

import (
	"sync"
)

// Internal events are events between States
// A window should never subscribe to an internal event
const (
	InternalCollectionSellected EventType = "internal:collection:selected"
	CollectionSelected          EventType = "collection:selected"
	CollectionChanged           EventType = "collection:changed"
	RequestSelected             EventType = "request:selected"
	RequestChanged              EventType = "request:changed"
	RequestExecuted             EventType = "request:executed"
	ResponseReceived            EventType = "response:received"
	ResponseError               EventType = "response:error"
)

type EventType string

type Event struct {
	Type EventType
	Data interface{}
}

type EventBus struct {
	asyncSubs  map[EventType][]chan Event
	subscribes map[EventType][]func(Event)
	mu         sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		asyncSubs:  make(map[EventType][]chan Event),
		subscribes: make(map[EventType][]func(Event)),
	}
}

func (b *EventBus) Subscribe(eventType EventType, f func(e Event)) {
	b.subscribes[eventType] = append(b.subscribes[eventType], f)
}

func (b *EventBus) Publish(event Event) {
	for _, f := range b.subscribes[event.Type] {
		f(event)
	}
}

func (b *EventBus) SubscribeAsync(eventType EventType) <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, 10)
	b.asyncSubs[eventType] = append(b.asyncSubs[eventType], ch)
	return ch
}

func (b *EventBus) PublishAsync(event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if chans, ok := b.asyncSubs[event.Type]; ok {
		for _, ch := range chans {
			go func(c chan Event) {
				c <- event
			}(ch)
		}
	}
}

