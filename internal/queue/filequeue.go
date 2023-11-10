package queue

import (
	"bytes"
	"encoding/gob"
	"errors"
	"sync"

	"github.com/andrescosta/workflew/internal/utils"
)

const (
	preffix = "qdata"
	suffix  = ".q"
	dir     = "data"
)

var queuesMap sync.Map

func GetFileBasedQueue[T any](id Id) *FileBasedQueue[T] {
	directory := queueDirectory(dir, id)
	queue, ok := queuesMap.Load(directory)
	if !ok {
		newQueue := &FileBasedQueue[T]{directory: directory}
		queue, _ = queuesMap.LoadOrStore(directory, newQueue)
	}
	return queue.(*FileBasedQueue[T])

}

type FileBasedQueue[T any] struct {
	directory string
	mutex     sync.Mutex
}

// aka Poor man queue
func NewDefaultFileBasedQueue[T any]() (*FileBasedQueue[T], error) {
	return &FileBasedQueue[T]{
		directory: dir,
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
	bdata, filename, err := utils.GetOldestFile(f.directory, preffix, suffix)
	if err != nil {
		var d T
		return d, errors.Join(errors.New("Error removing file"), err)
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
		utils.RenameFile(*filename, *filename+".error")
		return d, errors.Join(errors.New("Error encoding"), err)
	}
	utils.RemoveFile(*filename)
	return data, nil
}

func (f *FileBasedQueue[T]) writeData(data T) error {
	buffer := bytes.NewBuffer(make([]byte, 0))
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		return errors.Join(errors.New("Error encoding"), err)
	}
	err = utils.WriteToRandomFile(f.directory, preffix, suffix, buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func queueDirectory(directory string, id Id) string {
	return utils.BuildFullPath([]string{directory, id.Merchant, id.Name})
}
