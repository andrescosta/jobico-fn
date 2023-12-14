package repo

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/andrescosta/goico/pkg/ioutl"
	pb "github.com/andrescosta/jobico/api/types"
)

var (
	ErrFileExists = errors.New("file exists")
)

const (
	metFileExt = ".met"
	dirMeta    = "meta"
	dirFiles   = "content"
)

type FileRepo struct {
	dirFile string
	dirMeta string
}

func NewFileRepo(baseDir string) *FileRepo {
	return &FileRepo{
		dirFile: filepath.Join(baseDir, dirFiles),
		dirMeta: filepath.Join(baseDir, dirMeta),
	}
}

func (f *FileRepo) File(tenantId string, name string) ([]byte, error) {
	return file(name, f.dirFile, tenantId)
}

func file(name string, dirs ...string) ([]byte, error) {
	full := ioutl.BuildFullPath(dirs)
	res, err := os.ReadFile(ioutl.BuildPathWithFile(full, name))
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (f *FileRepo) Files() ([]*pb.TenantFiles, error) {
	dirs, err := ioutl.Dirs(f.dirFile)
	if err != nil {
		return nil, err
	}
	ts := make([]*pb.TenantFiles, 0)
	for _, dir := range dirs {
		fd := ioutl.BuildFullPath([]string{f.dirFile, dir.Name()})
		files, err := ioutl.Files(fd)
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
	if err := addFile(name, bytes, f.dirFile, tenantId); err != nil {
		return err
	}
	if err := f.WriteMetadataForFile(tenantId, name, fileType); err != nil {
		return err
	}
	return nil
}

func addFile(name string, bytes []byte, dirs ...string) error {
	full := ioutl.BuildFullPath(dirs)
	if err := ioutl.CreateDirsIfNotExist(full); err != nil {
		return err
	}
	fulPath := ioutl.BuildPathWithFile(full, name)
	e, err := ioutl.FileExists(fulPath)
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
	if err := addFile(name+metFileExt, buf.Bytes(), f.dirMeta, tenantId); err != nil {
		return err

	}
	return nil
}

func (f *FileRepo) GetMetadataForFile(tenantId string, name string) (*Metadata, error) {
	c, err := file(name+metFileExt, f.dirMeta, tenantId)
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
