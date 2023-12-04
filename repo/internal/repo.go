package repo

import (
	"os"

	"github.com/andrescosta/goico/pkg/iohelper"
	pb "github.com/andrescosta/jobico/api/types"
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

func (f *FileRepo) Files() ([]*pb.TenantFiles, error) {
	dirs, err := iohelper.GetDirs(f.Dir)
	if err != nil {
		return nil, err
	}
	ts := make([]*pb.TenantFiles, 0)
	for _, dir := range dirs {
		fd := iohelper.BuildFullPath([]string{f.Dir, dir.Name()})
		files, err := iohelper.GetFiles(fd)
		if err != nil {
			return nil, err
		}
		fs := make([]string, 0)
		for _, file := range files {
			fs = append(fs, file.Name())
		}
		ts = append(ts, &pb.TenantFiles{TenantId: dir.Name(), Files: fs})
	}
	return ts, nil
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
