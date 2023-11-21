package database

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/rs/zerolog"
	bolt "go.etcd.io/bbolt"
)

type Serializer[S any] interface {
	Serialize(uint64, S) ([]byte, error)
	Deserialize(uint64, []byte) (S, error)
}

type Table[S any] struct {
	db         *bolt.DB
	serializer Serializer[S]
	name       string
}

func Open[S any](ctx context.Context, path string, name string, serializer Serializer[S]) (*Table[S], error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	schema := &Table[S]{
		serializer: serializer,
		name:       name,
		db:         db,
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(schema.name))
		return err
	}); err != nil {
		return nil, err
	}

	return schema, nil
}

func (s *Table[S]) Add(_ context.Context, data S) (uint64, error) {
	var id uint64
	if err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.name))
		if b == nil {
			return fmt.Errorf("bucket does not exist")
		}
		var err error
		id, err = b.NextSequence()
		if err != nil {
			return err
		}
		buf, err := s.serializer.Serialize(id, data)
		if err != nil {
			return err
		}
		if err = b.Put(itob(id), buf); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Table[S]) Get(_ context.Context, id uint64) (S, error) {
	var data S
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.name))
		if b == nil {
			return fmt.Errorf("bucket does not exist")
		}
		var err error
		data, err = s.serializer.Deserialize(id, b.Get(itob(id)))
		if err != nil {
			return err
		}
		return nil
	})
	return data, nil
}
func (s *Table[S]) GetAll(ctx context.Context) ([]S, error) {
	var data []S
	data = make([]S, 0)
	//logger := zerolog.Ctx(ctx)
	if err := s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(s.name))
		if b == nil {
			return fmt.Errorf("bucket does not exist")
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			d, err := s.serializer.Deserialize(btoi(k), v)
			if err != nil {
				//			logger.Err(err)
				println(err)
			}
			data = append(data, d)
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return data, nil
}
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btoi(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}

func (s *Table[S]) Close(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	if err := s.db.Close(); err != nil {
		logger.Debug().Msgf("Error closing the DB:%s", err)
		return err
	}
	return nil
}
