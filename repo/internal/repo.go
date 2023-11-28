package repo

import (
	"os"

	"github.com/andrescosta/goico/pkg/iohelper"
)

type FileRepo struct {
	Dir string
}

func (f *FileRepo) File(tenantId string, name string) ([]byte, error) {
	dirs := iohelper.BuildFullPath([]string{f.Dir, tenantId})
	res, err := os.ReadFile(iohelper.BuildPathWithFile(dirs, name))
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (f *FileRepo) AddFile(tenantId string, name string, bytes []byte) error {
	dirs := iohelper.BuildFullPath([]string{f.Dir, tenantId})
	if err := iohelper.CreateDirIfNotExist(dirs); err != nil {
		return err
	}
	err := os.WriteFile(iohelper.BuildPathWithFile(dirs, name), bytes, os.ModeExclusive)
	if err != nil {
		return err
	} else {
		return nil
	}
}
