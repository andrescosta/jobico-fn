package repo

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/fs"
	"os"
	"syscall"

	"github.com/andrescosta/goico/pkg/iohelper"
	pb "github.com/andrescosta/jobico/api/types"
)

var (
	ErrFileExists = errors.New("file exists")
)

const (
	metFileExt = ".met"
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
		fs := make([]*pb.File, 0)
		for _, file := range files {
			m, err := f.GetMetadataForFile(dir.Name(), file.Name())
			if err != nil {
				return nil, err
			}

			fs = append(fs, &pb.File{Name: file.Name(), Type: pb.File_FileType(m.FileType)})
		}
		ts = append(ts, &pb.TenantFiles{TenantId: dir.Name(), Files: fs})
	}
	return ts, nil
}
func (f *FileRepo) AddFile(tenantId string, name string, fileType int32, bytes []byte) error {
	if err := f.addFile(tenantId, name, bytes); err != nil {
		return err
	}
	if err := f.WriteMetadataForFile(tenantId, name, fileType); err != nil {
		return err
	}
	return nil
}

func (f *FileRepo) addFile(tenantId string, name string, bytes []byte) error {
	dirs := iohelper.BuildFullPath([]string{f.Dir, tenantId})
	if err := iohelper.CreateDirIfNotExist(dirs); err != nil {
		return err
	}
	fulPath := iohelper.BuildPathWithFile(dirs, name)
	e, err := iohelper.FileExists(fulPath)
	if err != nil {
		return err
	}
	if e {
		return ErrFileExists
	}
	if err := os.WriteFile(fulPath, bytes, os.ModeExclusive); err != nil {
		return err
	}
	return nil
}

type Metadata struct {
	FileType int32
}

func (f *FileRepo) WriteMetadataForFile(tenantId string, name string, fileType int32) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(Metadata{FileType: fileType}); err != nil {
		return err
	}
	if err := f.addFile(tenantId, name+metFileExt, buf.Bytes()); err != nil {
		return err

	}
	return nil
}

func (f *FileRepo) GetMetadataForFile(tenant string, name string) (*Metadata, error) {
	c, err := f.File(tenant, name+metFileExt)
	if err != nil {
		pe, ok := err.(*fs.PathError)
		if ok && errors.Is(syscall.ERROR_FILE_NOT_FOUND, pe.Unwrap()) {
			return &Metadata{}, nil
		}
		return nil, err
	}
	buf := bytes.NewBuffer(c)
	dec := gob.NewDecoder(buf)
	metadata := Metadata{}
	if err := dec.Decode(&metadata); err != nil {
		return nil, err
	}
	return &metadata, nil
}
