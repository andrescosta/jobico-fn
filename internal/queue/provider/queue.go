package provider

import "errors"

var ErrQueueEmpty = errors.New("queue empty")

const MaxItems = 100

type Queue[T any] interface {
	Add(data T) error
	Remove() ([]T, error)
}
