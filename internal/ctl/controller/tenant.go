package controller

import (
	"github.com/andrescosta/goico/pkg/database"
	"github.com/andrescosta/goico/pkg/service/grpc/protoutil"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/ctl/dao"
	"google.golang.org/protobuf/proto"
)

const (
	tblTenant = "tenant"
)

type TenantController struct {
	daoCache *dao.DAOS
}

func NewTenantController(db *database.Database) *TenantController {
	return &TenantController{
		daoCache: dao.NewDAOS(db),
	}
}

func (c *TenantController) Close() error {
	return nil
}

func (c *TenantController) Tenants(in *pb.TenantsRequest) (*pb.TenantsReply, error) {
	if in.ID != nil {
		t, err := c.getTenant(*in.ID)
		if err != nil {
			return nil, err
		}
		if t != nil {
			return &pb.TenantsReply{Tenants: []*pb.Tenant{t}}, nil
		}
		return &pb.TenantsReply{}, nil
	}
	ts, err := c.getTenants()
	if err != nil {
		return nil, err
	}
	return &pb.TenantsReply{Tenants: ts}, nil
}

func (c *TenantController) AddTenant(in *pb.AddTenantRequest) (*pb.AddTenantReply, error) {
	mydao, err := c.daoCache.Generic(tblTenant, &pb.Tenant{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Tenant
	if err := mydao.Add(m); err != nil {
		return nil, err
	}
	return &pb.AddTenantReply{Tenant: in.Tenant}, nil
}

func (c *TenantController) getTenants() ([]*pb.Tenant, error) {
	mydao, err := c.daoCache.Generic(tblTenant, &pb.Tenant{})
	if err != nil {
		return nil, err
	}
	ms, err := mydao.All()
	if err != nil {
		return nil, err
	}
	tenants := protoutil.Slices[*pb.Tenant](ms)
	return tenants, nil
}

func (c *TenantController) getTenant(id string) (*pb.Tenant, error) {
	mydao, err := c.daoCache.Generic(tblTenant, &pb.Tenant{})
	if err != nil {
		return nil, err
	}
	ms, err := mydao.Get(id)
	if err != nil {
		return nil, err
	}
	if ms != nil {
		return (*ms).(*pb.Tenant), nil
	}
	return nil, nil
}
