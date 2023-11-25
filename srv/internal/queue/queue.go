package queue

type Id struct {
	QueueId    string
	MerchantId string
}

type Queue[T any] interface {
	Add(data T) error
	Remove() (T, error)
}

func GetQueue[T any](id Id) (Queue[T], error) {
	return GetFileBasedQueue[T](id), nil
}
