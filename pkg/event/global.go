package event

import "context"

var def = &Dispatcher{
	values: make(map[string]ListenerQueue),
}

func On(key string, listener Listener, priority Priority) {
	def.Register(key, ListenerPriority{
		Priority: priority,
		Listener: listener,
	})
}

func Register(name string, l ...ListenerPriority) {
	def.values[name] = ListenerQueue{
		listeners: append(def.values[name].listeners, l...),
	}
}

func AwaitFire(ctx context.Context, event Event) error {
	return def.AwaitFire(ctx, event)
}

func AsyncFire(ctx context.Context, event Event) {
	def.AsyncFire(ctx, event)
}

func Fire(ctx context.Context, event Event) error {
	return def.Fire(ctx, event)
}
