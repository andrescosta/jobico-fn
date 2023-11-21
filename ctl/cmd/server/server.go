package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/ctl/internal/database"
	"google.golang.org/protobuf/proto"
)

type Serializer struct{}

func (s *Serializer) Serialize(id uint64, q *types.QueueDef) ([]byte, error) {
	q.ID = strconv.FormatUint(id, 10)
	return proto.Marshal(q)
}

func (s *Serializer) Deserialize(id uint64, d []byte) (*types.QueueDef, error) {
	dd := types.QueueDef{}
	if err := proto.Unmarshal(d, &dd); err != nil {
		return nil, err
	}
	return &dd, nil
}

func main() {

	q := &types.QueueDef{
		QueueId: &types.QueueId{
			Name: "aaaa",
		},
	}

	s := Serializer{}
	schema, err := database.Open(context.Background(), ".\\db.db", "queue", &s)
	if err != nil {
		println(err)
	}
	i, err := schema.Add(context.Background(), q)
	if err != nil {
		println(err)
	}
	k, err := schema.Get(context.Background(), i)
	if err != nil {
		println(err)
	}
	println(k.String())
	ks, err := schema.GetAll(context.Background())
	if err != nil {
		println(err)
	}
	for _, kss := range ks {
		fmt.Println(kss)
	}
	schema.Close(context.Background())
}
