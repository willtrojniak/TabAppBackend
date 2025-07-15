package events

type EventDispatcher struct {
	listeners []func(Event)
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		listeners: make([]func(Event), 0),
	}
}

func Register[T Event](d *EventDispatcher, cb func(event T)) {
	d.listeners = append(d.listeners, func(e Event) {
		v, ok := e.(T)
		if ok {
			cb(v)
		}
	})
}

func Dispatch[T Event](d *EventDispatcher, event T) {
	for _, fn := range d.listeners {
		fn(event)
	}
}
