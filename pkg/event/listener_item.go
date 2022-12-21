package event

import "sort"

var _ sort.Interface = ByListenerPriority{}

type (
	ListenerPriority struct {
		Listener Listener
		Priority Priority
	}

	ByListenerPriority []ListenerPriority
)

func (b ByListenerPriority) Len() int {
	return len(b)
}
func (b ByListenerPriority) Less(i, j int) bool {
	return b[i].Priority.Value() > b[j].Priority.Value()
}
func (b ByListenerPriority) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
