package queue

import (
	"bytes"
	"encoding/gob"
	"errors"
	"sync"

	"github.com/andrescosta/goico/pkg/ioutil"
)

const (
	preffix = "qdata"

	suffix = ".q"

	DIR = "data"
)

var queuesMap sync.Map

type FileBasedQueue[T any] struct {
	directory string

	mutex sync.Mutex
}

func GetFileBasedQueue[T any](id string) *FileBasedQueue[T] {
	directory := queueDirectory(DIR, id)

	queue, ok := queuesMap.Load(directory)

	if !ok {
		newQueue := &FileBasedQueue[T]{directory: directory}

		queue, _ = queuesMap.LoadOrStore(directory, newQueue)
	}

	return queue.(*FileBasedQueue[T])
}

// aka Poor man queue

func NewDefaultFileBasedQueue[T any]() (*FileBasedQueue[T], error) {
	return &FileBasedQueue[T]{

		directory: DIR,
	}, nil
}

func (f *FileBasedQueue[T]) Add(data T) error {
	return f.writeData(data)
}

func (f *FileBasedQueue[T]) Remove() (T, error) {
	return f.readAndRemove()
}

func (f *FileBasedQueue[T]) readAndRemove() (T, error) {
	// sync the access to the "queue"

	f.mutex.Lock()

	defer f.mutex.Unlock()

	bdata, filename, err := ioutil.OldestFile(f.directory, preffix, suffix)

	if err != nil {
		var d T

		return d, errors.Join(errors.New("error removing file"), err)
	}

	if bdata == nil {
		var d T

		return d, nil
	}

	buffer := bytes.NewBuffer(bdata)

	decoder := gob.NewDecoder(buffer)

	var data T

	if err = decoder.Decode(&data); err != nil {
		var d T

		if errR := ioutil.Rename(*filename, *filename+".error"); errR != nil {
			err = errors.Join(errR, err)
		}

		return d, errors.Join(errors.New("error decoding"), err)
	}

	if err := ioutil.Remove(*filename); err != nil {
		var t T
		return t, err
	}

	return data, nil
}

func (f *FileBasedQueue[T]) writeData(data T) error {
	buffer := bytes.NewBuffer(make([]byte, 0))

	encoder := gob.NewEncoder(buffer)

	err := encoder.Encode(data)

	if err != nil {
		return errors.Join(errors.New("error encoding"), err)
	}

	err = ioutil.WriteToRandomFile(f.directory, preffix, suffix, buffer.Bytes())

	if err != nil {
		return err
	}

	return nil
}

func queueDirectory(directory string, id string) string {
	return ioutil.BuildFullPath([]string{directory, id})
}
