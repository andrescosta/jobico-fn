package dao

import (
	"github.com/andrescosta/goico/pkg/database"
	"google.golang.org/protobuf/proto"
)

type DAO[T proto.Message] struct {
	table *database.Table[T]
}

func New[T proto.Message](db *database.Database, tableName string, tenant string, m database.Marshaler[T]) *DAO[T] {
	t := database.NewTable(db, tableName, tenant, m)
	return &DAO[T]{
		table: t,
	}
}

func (q *DAO[T]) All() ([]T, error) {
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
