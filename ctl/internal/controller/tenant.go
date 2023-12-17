package controller

import (
	"context"

	"github.com/andrescosta/goico/pkg/convert"
	"github.com/andrescosta/goico/pkg/database"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/ctl/internal/dao"
	"google.golang.org/protobuf/proto"
)

const (
	tblTenant = "tenant"
)

type Tenant struct {
	daoCache *dao.Cache
}

func NewTenant(db *database.Database) *Tenant {
	return &Tenant{
		daoCache: dao.NewCache(db),
	}
}

func (c *Tenant) Close() {
}

func (c *Tenant) GetTenants(ctx context.Context, in *pb.GetTenantsRequest) (*pb.GetTenantsReply, error) {
	if in.ID != nil {
		t, err := c.getTenant(ctx, *in.ID)
		if err != nil {
			return nil, err
		}
		if t != nil {
			return &pb.GetTenantsReply{Tenants: []*pb.Tenant{t}}, nil
		}
		return &pb.GetTenantsReply{}, nil
	}
	ts, err := c.getTenants(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetTenantsReply{Tenants: ts}, nil
}
func (c *Tenant) AddTenant(ctx context.Context, in *pb.AddTenantRequest) (*pb.AddTenantReply, error) {
	mydao, err := c.daoCache.GetGeneric(ctx, tblTenant, &pb.Tenant{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Tenant
	_, err = mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	return &pb.AddTenantReply{Tenant: in.Tenant}, nil
}
func (c *Tenant) getTenants(ctx context.Context) ([]*pb.Tenant, error) {
	mydao, err := c.daoCache.GetGeneric(ctx, tblTenant, &pb.Tenant{})
	if err != nil {
		return nil, err
	}
	ms, err := mydao.All(ctx)
	if err != nil {
		return nil, err
	}
	tenants := convert.Slices[proto.Message, *pb.Tenant](ms)
	return tenants, nil
}
func (c *Tenant) getTenant(ctx context.Context, id string) (*pb.Tenant, error) {
	mydao, err := c.daoCache.GetGeneric(ctx, tblTenant, &pb.Tenant{})
	if err != nil {
		return nil, err
	}
	ms, err := mydao.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if ms != nil {
		return (*ms).(*pb.Tenant), nil
	}
	return nil, nil
}
