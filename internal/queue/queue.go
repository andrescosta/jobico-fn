package queue

type Queue[T any] interface {
	Add(data T) error
	Remove() (T, error)
}

func GetDefault[T any]() Queue[T] {
	return NewDefaultFileBasedQueue[T]()
}
