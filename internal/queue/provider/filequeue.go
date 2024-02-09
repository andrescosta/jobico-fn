package provider

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"sync"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/ioutil"
)

const (
	preffix = "qdata"
	suffix  = ".q"
	dir     = "data"
)

type FileQueue[T any] struct {
	directory string
	mutex     sync.Mutex
}

func NewFileQueue[T any](id string) *FileQueue[T] {
	directory := queueDirectory(dir, id)
	return &FileQueue[T]{directory: directory}
}

func (f *FileQueue[T]) Add(data T) error {
	return f.writeData(data)
}

func (f *FileQueue[T]) Remove() ([]T, error) {
	return f.readAndRemove()
}

func (f *FileQueue[T]) readAndRemove() ([]T, error) {
	// sync the access to the "queue"
	f.mutex.Lock()
	defer f.mutex.Unlock()
	files, err := ioutil.ReadOldestFiles(f.directory, preffix, suffix, MaxItems)
	if err != nil {
		var d []T
		return d, errors.Join(errors.New("error removing file"), err)
	}
	if len(files) == 0 {
		var d []T
		return d, ErrQueueEmpty
	}
	ts := make([]T, len(files))
	for i, f := range files {
		bdata := f.Bytes
		filename := f.Name
		buffer := bytes.NewBuffer(bdata)
		decoder := gob.NewDecoder(buffer)
		var data T
		if err := decoder.Decode(&data); err != nil {
			var d []T
			// if the file cannot be decoded, we rename it to [file].error for further processing
			if errR := os.Rename(filename, filename+".error"); errR != nil {
				err = errors.Join(errR, err)
			}
			return d, errors.Join(errors.New("error decoding"), err)
		}
		if err := os.Remove(filename); err != nil {
			var t []T
			return t, err
		}
		ts[i] = data
	}
	return ts, nil
}

func (f *FileQueue[T]) writeData(data T) error {
	buffer := bytes.NewBuffer(make([]byte, 0))
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(data); err != nil {
		return errors.Join(errors.New("error encoding"), err)
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()
	if _, err := ioutil.WriteToRandomFile(f.directory, preffix, suffix, buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

func queueDirectory(directory string, id string) string {
	return env.WorkdirPlus(directory, id)
}
