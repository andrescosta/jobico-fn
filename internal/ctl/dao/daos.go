package dao

import (
	"sync"

	"github.com/andrescosta/goico/pkg/database"
	"google.golang.org/protobuf/proto"
)

type DAOS struct {
	mu   *sync.Mutex
	daos map[string]*DAO[proto.Message]
	db   *database.Database
}

func NewDAOS(db *database.Database) *DAOS {
	return &DAOS{
		daos: make(map[string]*DAO[proto.Message]),
		db:   db,
		mu:   &sync.Mutex{},
	}
}

func (c *DAOS) Generic(entity string, message proto.Message) (*DAO[proto.Message], error) {
	return c.ForTenant(entity, entity, message)
}

func (c *DAOS) ForTenant(tenant string, entity string, message proto.Message) (*DAO[proto.Message], error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	dao, ok := c.daos[tenant]
	if !ok {
		dao = New(c.db, entity, tenant,
			&ProtoMessageMarshaller{
				prototype: message,
			})
		c.daos[tenant] = dao
	}
	return dao, nil
}
