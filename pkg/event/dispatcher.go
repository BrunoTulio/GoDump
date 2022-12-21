package event

import "context"

type Dispatcher struct {
	values map[string]ListenerQueue
}

func (e *Dispatcher) On(key string, listener Listener, priority Priority) {
	e.Register(key, ListenerPriority{
		Priority: priority,
		Listener: listener,
	})
}

func (e *Dispatcher) Register(name string, l ...ListenerPriority) {
	e.values[name] = ListenerQueue{
		listeners: append(e.values[name].listeners, l...),
	}
}

// AwaitFire async fire event by 'go' keywords, but will wait return result
func (e *Dispatcher) AwaitFire(ctx context.Context, event Event) (err error) {
	ch := make(chan error)

	go func(ctx context.Context, event Event) {
		err := e.Fire(ctx, event)
		ch <- err
	}(ctx, event)

	err = <-ch
	close(ch)
	return
}

func (e *Dispatcher) AsyncFire(ctx context.Context, event Event) {
	go func(ctx context.Context, event Event) {
		_ = e.Fire(ctx, event)
	}(ctx, event)
}

func (e *Dispatcher) Fire(ctx context.Context, event Event) error {
	if queue, ok := e.values[event.Key()]; ok && !queue.IsEmpty() {
		for _, item := range queue.Sort().Items() {
			if err := item.Listener.Handler(ctx, event); err != nil {
				return err
			}
		}
	}
	return nil
}
