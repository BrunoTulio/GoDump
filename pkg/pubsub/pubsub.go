package pubsub

import "sync"

type (
	PubSub interface {
		Subscribe(topic string) <-chan any
		Publish(topic string, value interface{})
		Close()
	}

	pubSub struct {
		mu     sync.Mutex
		subs   map[string][]chan any
		quit   chan struct{}
		closed bool
	}
)

// Close implements SSEListener.
func (s *pubSub) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	s.closed = true
	close(s.quit)

	for _, ch := range s.subs {
		for _, sub := range ch {
			close(sub)
		}
	}
}

// Publish implements SSEListener.
func (s *pubSub) Publish(topic string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	for _, ch := range s.subs[topic] {
		ch <- value
	}

}

// Subscribe implements SSEListener.
func (s *pubSub) Subscribe(topic string) <-chan any {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	ch := make(chan any)
	s.subs[topic] = append(s.subs[topic], ch)
	return ch

}

func NewPubSub() PubSub {
	return &pubSub{
		subs: make(map[string][]chan any),
		quit: make(chan struct{}),
	}
}
