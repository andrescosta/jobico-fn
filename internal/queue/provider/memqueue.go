package provider

import (
	"github.com/andrescosta/goico/pkg/collection"
)

type MemBasedQueue[T any] struct {
	q *collection.SyncQueue[T]
}

func NewMemBasedQueue[T any]() (*MemBasedQueue[T], error) {
	return &MemBasedQueue[T]{
		q: collection.NewSyncQueue[T](),
	}, nil
}

func (f *MemBasedQueue[T]) Add(data T) error {
	f.q.Queue(data)
	return nil
}

func (f *MemBasedQueue[T]) Remove() ([]T, error) {
	if f.q.Size() == 0 {
		var t []T
		return t, ErrQueueEmpty
	}
	return f.q.DequeueN(MaxItems), nil
}
