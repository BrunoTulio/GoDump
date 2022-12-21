package event

import "context"

type Listener interface {
	Handler(ctx context.Context, e Event) error
}
