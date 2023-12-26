package dao

import (
	"context"

	"github.com/andrescosta/goico/pkg/database"
	"google.golang.org/protobuf/proto"
)

type Cache struct {
	daos map[string]*DAO[proto.Message]
	db   *database.Database
}

func NewCache(db *database.Database) *Cache {
	return &Cache{
		daos: make(map[string]*DAO[proto.Message]),
		db:   db,
	}
}

func (c *Cache) GetGeneric(ctx context.Context, entity string, message proto.Message) (*DAO[proto.Message], error) {
	return c.GetForTenant(ctx, entity, entity, message)
}

func (c *Cache) GetForTenant(ctx context.Context, tenant string, entity string, message proto.Message) (*DAO[proto.Message], error) {
	mydao, ok := c.daos[tenant]
	if !ok {
		var err error
		mydao, err = NewDAO(ctx, c.db, tenant+"/"+entity,
			&ProtoMessageMarshaler{
				prototype: message,
			})
		if err != nil {
			return nil, err
		}
		c.daos[tenant] = mydao
	}
	return mydao, nil
}
