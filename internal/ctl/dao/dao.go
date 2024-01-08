package dao

import (
	"context"

	"github.com/andrescosta/goico/pkg/database"
	"google.golang.org/protobuf/proto"
)

type DAO[T proto.Message] struct {
	table *database.Table[T]
}

func NewDAO[T proto.Message](db *database.Database, tableName string, m database.Marshaler[T]) (*DAO[T], error) {
	table, err := database.CreateTableIfNotExist(db, tableName, m)
	if err != nil {
		return nil, err
	}
	res := DAO[T]{
		table: table,
	}
	return &res, nil
}

func (q *DAO[T]) All(ctx context.Context) ([]T, error) {
	return q.table.All()
}

func (q *DAO[T]) Get(id string) (*T, error) {
	return q.table.Get(id)
}

func (q *DAO[T]) Add(data T) error {
	return q.table.Add(data)
}

func (q *DAO[T]) Update(data T) error {
	return q.table.Update(data)
}

func (q *DAO[T]) Delete(id string) error {
	return q.table.Delete(id)
}
