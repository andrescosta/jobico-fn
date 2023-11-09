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
)

type FileBasedQueue[T any] struct {
	directory string
	mutex     sync.Mutex
}

func NewDefaultFileBasedQueue[T any]() *FileBasedQueue[T] {
	return &FileBasedQueue[T]{
		directory: ".",
	}
}

func (f *FileBasedQueue[T]) Add(data T) error {
	return f.writeData(data)
}

func (f *FileBasedQueue[T]) Remove() (T, error) {
	return f.readAndRemove()
}

func (f *FileBasedQueue[T]) readAndRemove() (T, error) {
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
