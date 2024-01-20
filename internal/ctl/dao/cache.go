package dao

import (
	"sync"

	"github.com/andrescosta/goico/pkg/database"
	"google.golang.org/protobuf/proto"
)

type Cache struct {
	mu   *sync.Mutex
	daos map[string]*DAO[proto.Message]
	db   *database.Database
}

func NewCache(db *database.Database) *Cache {
	return &Cache{
		daos: make(map[string]*DAO[proto.Message]),
		db:   db,
		mu:   &sync.Mutex{},
	}
}

func (c *Cache) GetGeneric(entity string, message proto.Message) (*DAO[proto.Message], error) {
	return c.GetForTenant(entity, entity, message)
}

func (c *Cache) GetForTenant(tenant string, entity string, message proto.Message) (*DAO[proto.Message], error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	dao, ok := c.daos[tenant]
	if !ok {
		dao = NewDAO(c.db, entity, tenant,
			&ProtoMessageMarshaller{
				prototype: message,
			})
		c.daos[tenant] = dao
	}
	return dao, nil
}
