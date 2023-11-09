package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func WriteToRandomFile(directory, preffix, suffix string, data []byte) error {
	f := directory + string(os.PathSeparator) + randomFileName(preffix, suffix)
	return WriteToFile(f, data)
}

func WriteToFile(file string, data []byte) error {
	f, err := os.Create(file)
	if err != nil {
		return errors.Join(errors.New("Error creating file"), err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			// log
		}
	}()

	w := bufio.NewWriter(f)

	if _, err := w.Write(data); err != nil {
		return errors.Join(errors.New("Error writing file"), err)
	}

	if err := w.Flush(); err != nil {
		return errors.Join(errors.New("Error flushing data"), err)
	}
	return nil
}

func getFilesSorted(directory, preffix, suffix string) ([]os.DirEntry, error) {
	var files []fs.DirEntry
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if strings.HasPrefix(d.Name(), preffix) && strings.HasSuffix(d.Name(), suffix) {
			files = append(files, d)
		}
		return nil
	})

	if err != nil {
		return nil, errors.Join(errors.New("Error getting files"), err)
	}
	sort.Slice(files, func(i, j int) bool {
		in1, _ := files[i].Info()
		in2, _ := files[j].Info()
		return in1.ModTime().Unix() < in2.ModTime().Unix()
	})
	return files, nil
}

func getSomeFilesSorted(directory, preffix, suffix string, n int) ([]os.DirEntry, error) {
	files, err := getFilesSorted(directory, preffix, suffix)
	if err != nil {
		return nil, err
	}
	if len(files) < n {
		return files, nil
	} else {
		return files[:n], nil
	}
}

func GetOldestFile(directory, preffix, suffix string) ([]byte, *string, error) {
	files, err := getSomeFilesSorted(directory, preffix, suffix, 1)
	if err != nil {
		return nil, nil, err
	}
	if len(files) == 0 {
		return nil, nil, nil
	}
	filename := directory + string(os.PathSeparator) + files[0].Name()
	data, err := ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	return data, &filename, nil
}

func ReadFile(file string) ([]byte, error) {
	return os.ReadFile(file)
}

func RemoveFile(file string) error {
	return os.Remove(file)
}

func randomFileName(preffix, suffix string) string {
	s := strconv.FormatUint(uint64(random.Uint32()), 10)
	return preffix + s + suffix
}
