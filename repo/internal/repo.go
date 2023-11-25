package repo

import (
	"os"

	"github.com/andrescosta/goico/pkg/io"
)

type FileRepo struct {
	Dir string
}

func (f *FileRepo) File(tenantId string, name string) ([]byte, error) {
	dirs := io.BuildFullPath([]string{f.Dir, tenantId})
	res, err := os.ReadFile(io.BuildPathWithFile(dirs, name))
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (f *FileRepo) AddFile(tenantId string, name string, bytes []byte) error {
	dirs := io.BuildFullPath([]string{f.Dir, tenantId})
	if err := io.CreateDirIfNotExist(dirs); err != nil {
		return err
	}
	err := os.WriteFile(io.BuildPathWithFile(dirs, name), bytes, os.ModeExclusive)
	if err != nil {
		return err
	} else {
		return nil
	}
}
