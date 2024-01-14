package dao

import (
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

func (c *Cache) GetGeneric(entity string, message proto.Message) (*DAO[proto.Message], error) {
	return c.GetForTenant(entity, entity, message)
}

func (c *Cache) GetForTenant(tenant string, entity string, message proto.Message) (*DAO[proto.Message], error) {
	mydao, ok := c.daos[tenant]
	if !ok {
		var err error
		mydao, err = NewDAO(c.db, tenant+"/"+entity,
			&ProtoMessageMarshaller{
				prototype: message,
			})
		if err != nil {
			return nil, err
		}
		c.daos[tenant] = mydao
	}
	return mydao, nil
}
