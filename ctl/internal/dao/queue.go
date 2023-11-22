package dao

import (
	"context"
	"strconv"

	"github.com/andrescosta/goico/pkg/database"
	"github.com/andrescosta/workflew/api/types"
	"google.golang.org/protobuf/proto"
)

type QueueDAO struct {
	table *database.Table[*types.QueueDef]
}

func NewQueueDAO(ctx context.Context, path string) (*QueueDAO, error) {
	table, err := database.Open(context.Background(), ".\\db.db", "queue", &QueueDefSerializer{})
	if err != nil {
		return nil, err
	}
	return &QueueDAO{
		table: table,
	}, nil
}

func (q *QueueDAO) Close(ctx context.Context) error {
	return q.table.Close(ctx)
}

func (q *QueueDAO) GetQueueDefs(ctx context.Context) ([]*types.QueueDef, error) {
	return q.table.GetAll(ctx)
}

func (q *QueueDAO) GetQueueDef(ctx context.Context, id uint64) (*types.QueueDef, error) {
	return q.table.Get(ctx, id)
}

func (q *QueueDAO) AddQueueDef(ctx context.Context, data *types.QueueDef) (uint64, error) {
	return q.table.Add(ctx, data)
}

type QueueDefSerializer struct{}

func (s *QueueDefSerializer) Serialize(id uint64, q *types.QueueDef) ([]byte, error) {
	q.ID = strconv.FormatUint(id, 10)
	return proto.Marshal(q)
}

func (s *QueueDefSerializer) Deserialize(id uint64, d []byte) (*types.QueueDef, error) {
	dd := types.QueueDef{}
	if err := proto.Unmarshal(d, &dd); err != nil {
		return nil, err
	}
	return &dd, nil
}
