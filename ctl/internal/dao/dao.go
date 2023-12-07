package dao

import (
	"context"

	"github.com/andrescosta/goico/pkg/database"
	"google.golang.org/protobuf/proto"
)

type DAO[T proto.Message] struct {
	table *database.Table[T]
}

func NewDAO[T proto.Message](ctx context.Context, db *database.Database, tableName string, m database.Marshaler[T]) (*DAO[T], error) {
	table, err := database.GetTable(ctx, db, tableName, m)
	if err != nil {
		return nil, err
	}
	var r DAO[T] = DAO[T]{
		table: table,
	}
	return &r, nil
}

func (q *DAO[T]) All(ctx context.Context) ([]T, error) {
	return q.table.All(ctx)
}

func (q *DAO[T]) Get(ctx context.Context, id string) (*T, error) {
	return q.table.Get(ctx, id)
}

func (q *DAO[T]) Add(ctx context.Context, data T) (uint64, error) {
	i, err := q.table.Add(ctx, data)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (q *DAO[T]) Update(ctx context.Context, data T) error {
	if err := q.table.Update(ctx, data); err != nil {
		return err
	}
	return nil
}
