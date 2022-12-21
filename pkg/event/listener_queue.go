package event

import (
	"fmt"
	"sort"
)

type ListenerQueue struct {
	listeners []ListenerPriority
}

func (l *ListenerQueue) Len() int {
	return len(l.listeners)
}

func (l *ListenerQueue) IsEmpty() bool {
	return len(l.listeners) == 0
}

func (l *ListenerQueue) Push(li ListenerPriority) *ListenerQueue {
	l.listeners = append(l.listeners, li)
	return l
}

func (l *ListenerQueue) Remove(li Listener) {
	if li == nil {
		return
	}

	ptrVal := fmt.Sprintf("%p", li)

	var newItems []ListenerPriority
	for _, li := range l.listeners {
		liPtrVal := fmt.Sprintf("%p", li.Listener)
		if liPtrVal == ptrVal {
			continue
		}

		newItems = append(newItems, li)
	}

	l.listeners = newItems
}

func (l *ListenerQueue) Sort() *ListenerQueue {
	ls := ByListenerPriority(l.listeners)
	if !sort.IsSorted(ls) {
		sort.Sort(ls)
	}
	return l
}

func (l ListenerQueue) Items() []ListenerPriority {
	return l.listeners
}

func (lq *ListenerQueue) Clear() {
	lq.listeners = lq.listeners[:0]
}
