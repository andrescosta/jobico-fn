package dao

import (
	"context"

	"github.com/andrescosta/goico/pkg/database"
)

type DAO[T any] struct {
	table *database.Table[T]
}

func NewDAO[T any](ctx context.Context, db *database.Database, tableName string, m database.Marshaler[T]) (*DAO[T], error) {
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

func (q *DAO[T]) Get(ctx context.Context, id uint64) (T, error) {
	return q.table.Get(ctx, id)
}

func (q *DAO[T]) Add(ctx context.Context, data T) (uint64, error) {
	return q.table.Add(ctx, data)
}
