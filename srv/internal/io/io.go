package io

import (
	"bufio"
	"errors"
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

func WriteToRandomFile(path, preffix, suffix string, data []byte) error {
	CreateDirIfNotExist(path)
	fullpath := BuildPathWithFile(path, randomFileName(preffix, suffix))
	return WriteToFile(fullpath, data)
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

/*func GetDirectoriesEntry(path os.DirEntry) ([]os.DirEntry, error) {
return GetDirectories(
}*/

func GetSubDirectories(path string) ([]string, error) {
	dires, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	dirs := make([]string, 0)
	for _, d := range dires {
		dirs = append(dirs, d.Name())
	}
	return dirs, nil
}

func getFilesSorted(path string, preffix, suffix string) ([]os.DirEntry, error) {
	var files []fs.DirEntry

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return files, nil
	}

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
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

func getNFilesSorted(path, preffix, suffix string, n int) ([]os.DirEntry, error) {
	files, err := getFilesSorted(path, preffix, suffix)
	if err != nil {
		return nil, err
	}
	if len(files) < n {
		return files, nil
	} else {
		return files[:n], nil
	}
}

func GetOldestFile(path, preffix, suffix string) ([]byte, *string, error) {
	files, err := getNFilesSorted(path, preffix, suffix, 1)
	if err != nil {
		return nil, nil, err
	}
	if len(files) == 0 {
		return nil, nil, nil
	}
	filename := BuildFullPathWithFile([]string{path}, files[0].Name())
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

func RenameFile(file, newFile string) error {
	return os.Rename(file, newFile)
}

func randomFileName(preffix, suffix string) string {
	s := strconv.FormatUint(uint64(random.Uint32()), 10)
	return preffix + s + suffix
}

func CreateDirIfNotExist(directory string) error {
	return os.MkdirAll(directory, os.ModeExclusive)
}

func BuildFullPath(directories []string) string {
	return filepath.Join(directories...)
}

func BuildFullPathWithFile(directories []string, file string) string {
	return filepath.Join(BuildFullPath(directories), file)
}

func BuildPathWithFile(path, file string) string {
	return filepath.Join(path, file)
}
