package provider

import (
	"errors"
	"sync"

	pb "github.com/andrescosta/jobico/internal/api/types"
)

type MemRepo struct {
	mutex   *sync.RWMutex
	mapFile map[string]map[string][]byte
	mapMeta map[string]map[string]*Metadata
}

func NewMemRepo() *MemRepo {
	return &MemRepo{
		mapFile: make(map[string]map[string][]byte),
		mapMeta: make(map[string]map[string]*Metadata),
		mutex:   &sync.RWMutex{},
	}
}

func (m *MemRepo) Add(tenant string, name string, fileType int32, bytes []byte) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	td, ok := m.mapFile[tenant]
	if !ok {
		td = make(map[string][]byte)
		m.mapFile[tenant] = td
	}
	_, ok = td[name]
	if !ok {
		td[name] = bytes
	}
	tm, ok := m.mapMeta[tenant]
	if !ok {
		tm = make(map[string]*Metadata)
		m.mapMeta[tenant] = tm
	}
	_, ok = tm[name]
	if !ok {
		tm[name] = &Metadata{FileType: fileType}
	}
	return nil
}

var ErrNotFound = errors.New("not found")

func (m *MemRepo) File(tenant string, name string) ([]byte, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	td, ok := m.mapFile[tenant]
	if !ok {
		return nil, ErrNotFound
	}
	tf, ok := td[name]
	if !ok {
		return nil, ErrNotFound
	}
	return tf, nil
}

func (m *MemRepo) GetMetadataForFile(tenant string, name string) (*Metadata, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	td, ok := m.mapMeta[tenant]
	if !ok {
		return nil, ErrNotFound
	}
	tm, ok := td[name]
	if !ok {
		return nil, ErrNotFound
	}
	return tm, nil
}

func (m *MemRepo) Files() ([]*pb.TenantFiles, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	ts := make([]*pb.TenantFiles, 0)
	for kt, vf := range m.mapMeta {
		fs := make([]*pb.File, 0)
		for kf, vm := range vf {
			fs = append(fs, &pb.File{Name: kf, Type: pb.File_FileType(vm.FileType)})
		}
		ts = append(ts, &pb.TenantFiles{Tenant: kt, Files: fs})
	}
	return ts, nil
}
