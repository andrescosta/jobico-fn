package provider

import "errors"

var ErrQueueEmpty = errors.New("queue empty")

type Queue[T any] interface {
	Add(data T) error
	Remove() (T, error)
}
