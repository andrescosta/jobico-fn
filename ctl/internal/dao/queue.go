package dao

import (
	"context"

	"github.com/andrescosta/goico/pkg/database"
	"google.golang.org/protobuf/proto"
)

type newProtoMessage func() proto.Message

type ProtoMessageDAO struct {
	table *database.Table[*proto.Message]
}

type ProtoMessageSerializer struct {
	newMessage newProtoMessage
}

func NewProtoMessageDAO(ctx context.Context, db *database.Database, tableName string, message proto.Message) (*ProtoMessageDAO, error) {
	table, err := database.GetTable(ctx, db, tableName, &ProtoMessageSerializer{
		newMessage: func() proto.Message {
			return message
		},
	})
	if err != nil {
		return nil, err
	}
	return &ProtoMessageDAO{
		table: table,
	}, nil
}

func (q *ProtoMessageDAO) All(ctx context.Context) ([]*proto.Message, error) {
	return q.table.All(ctx)
}

func (q *ProtoMessageDAO) Get(ctx context.Context, id uint64) (*proto.Message, error) {
	return q.table.Get(ctx, id)
}

func (q *ProtoMessageDAO) Add(ctx context.Context, data *proto.Message) (uint64, error) {
	return q.table.Add(ctx, data)
}

func (s *ProtoMessageSerializer) Serialize(id uint64, q *proto.Message) ([]byte, error) {
	//q.ID = strconv.FormatUint(id, 10)
	//reflect.ValueOf(q).FieldByName("ID").SetString(strconv.FormatUint(id, 10))
	return proto.Marshal(*q)
}

func (s *ProtoMessageSerializer) Deserialize(id uint64, d []byte) (*proto.Message, error) {
	dd := s.newMessage()
	if err := proto.Unmarshal(d, dd); err != nil {
		return nil, err
	}
	return &dd, nil
}
