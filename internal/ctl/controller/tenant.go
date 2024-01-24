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
	daoCache *dao.Cache
}

func NewTenantController(db *database.Database) *TenantController {
	return &TenantController{
		daoCache: dao.NewCache(db),
	}
}

func (c *TenantController) Close() {
}

func (c *TenantController) GetTenants(in *pb.GetTenantsRequest) (*pb.GetTenantsReply, error) {
	if in.ID != nil {
		t, err := c.getTenant(*in.ID)
		if err != nil {
			return nil, err
		}
		if t != nil {
			return &pb.GetTenantsReply{Tenants: []*pb.Tenant{t}}, nil
		}
		return &pb.GetTenantsReply{}, nil
	}
	ts, err := c.getTenants()
	if err != nil {
		return nil, err
	}
	return &pb.GetTenantsReply{Tenants: ts}, nil
}

func (c *TenantController) AddTenant(in *pb.AddTenantRequest) (*pb.AddTenantReply, error) {
	mydao, err := c.daoCache.GetGeneric(tblTenant, &pb.Tenant{})
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
	mydao, err := c.daoCache.GetGeneric(tblTenant, &pb.Tenant{})
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
	mydao, err := c.daoCache.GetGeneric(tblTenant, &pb.Tenant{})
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
