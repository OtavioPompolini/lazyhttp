package state

import (
	"sync"
)

// Internal events are events between States
// A window should never subscribe to an internal event
const (
	InternalCollectionSelected EventType = "internal:collection:selected"
	CollectionSelected         EventType = "collection:selected"
	CollectionChanged          EventType = "collection:changed"
	RequestSelected            EventType = "request:selected"
	RequestChanged             EventType = "request:changed"
	RequestExecuted            EventType = "request:executed"
	ResponseReceived           EventType = "response:received"
	ResponseError              EventType = "response:error"
	AlertMessage               EventType = "alert:message"
)

type EventType string

type Event struct {
	Type EventType
	Data interface{}
}

type EventBus struct {
	asyncHandlers map[EventType][]chan Event
	handlers      map[EventType][]func(Event)
	mu            sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		asyncHandlers: make(map[EventType][]chan Event),
		handlers:      make(map[EventType][]func(Event)),
	}
}

func (b *EventBus) Subscribe(eventType EventType, f func(e Event)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], f)
}

func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	handlers := make([]func(Event), len(b.handlers[event.Type]))
	copy(handlers, b.handlers[event.Type])
	b.mu.RUnlock()

	for _, f := range handlers {
		f(event)
	}
}

func (b *EventBus) SubscribeAsync(eventType EventType) <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, 10)
	b.asyncHandlers[eventType] = append(b.asyncHandlers[eventType], ch)
	return ch
}

func (b *EventBus) PublishAsync(event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if chans, ok := b.asyncHandlers[event.Type]; ok {
		for _, ch := range chans {
			go func(c chan Event) {
				c <- event
			}(ch)
		}
	}
}
