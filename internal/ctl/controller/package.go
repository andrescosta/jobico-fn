package controller

import (
	"context"

	"github.com/andrescosta/goico/pkg/database"
	"github.com/andrescosta/goico/pkg/service/grpc/protoutil"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/ctl/dao"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/protobuf/proto"
)

const (
	tblPackage = "package"
)

type PackageController struct {
	ctx              context.Context
	daoCache         *dao.DAOS
	bJobPackage      *grpchelper.GrpcBroadcaster[*pb.UpdateToPackagesStrReply, proto.Message]
	tenantController *TenantController
}

func NewPackageController(ctx context.Context, db *database.Database) *PackageController {
	return &PackageController{
		ctx:              ctx,
		daoCache:         dao.NewDAOS(db),
		bJobPackage:      grpchelper.StartBroadcaster[*pb.UpdateToPackagesStrReply, proto.Message](ctx),
		tenantController: NewTenantController(db),
	}
}

func (c *PackageController) Close() error {
	return c.bJobPackage.Stop()
}

func (c *PackageController) GetPackages(in *pb.PackagesRequest) (*pb.PackagesReply, error) {
	if in.ID != nil {
		p, err := c.getPackage(in.Tenant, *in.ID)
		if err != nil {
			return nil, err
		}
		if p != nil {
			return &pb.PackagesReply{Packages: []*pb.JobPackage{p}}, nil
		}
		return &pb.PackagesReply{}, nil
	}
	packages, err := c.getPackages(in.Tenant)
	if err != nil {
		return nil, err
	}
	return &pb.PackagesReply{Packages: packages}, nil
}

func (c *PackageController) GetAllPackages() (*pb.AllPackagesReply, error) {
	ms, err := c.tenantController.getTenants()
	if err != nil {
		return nil, err
	}
	packages := make([]*pb.JobPackage, 0)
	for _, me := range ms {
		mydao, err := c.daoCache.ForTenant(me.ID, tblPackage, &pb.JobPackage{})
		if err != nil {
			return nil, err
		}
		ms, err := mydao.All()
		if err != nil {
			return nil, err
		}
		ps := protoutil.Slices[*pb.JobPackage](ms)
		packages = append(packages, ps...)
	}
	return &pb.AllPackagesReply{Packages: packages}, nil
}

func (c *PackageController) AddPackage(ctx context.Context, in *pb.AddPackageRequest) (*pb.AddPackageReply, error) {
	mydao, err := c.daoCache.ForTenant(in.Package.Tenant, tblPackage, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Package
	if err := mydao.Add(m); err != nil {
		return nil, err
	}
	c.broadcastAdd(ctx, in.Package)
	return &pb.AddPackageReply{Package: in.Package}, nil
}

func (c *PackageController) UpdatePackage(ctx context.Context, in *pb.UpdatePackageRequest) (*pb.Void, error) {
	mydao, err := c.daoCache.ForTenant(in.Package.Tenant, tblPackage, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Package
	err = mydao.Update(m)
	if err != nil {
		return nil, err
	}
	c.broadcastUpdate(ctx, in.Package)
	return &pb.Void{}, nil
}

func (c *PackageController) DeletePackage(ctx context.Context, in *pb.DeletePackageRequest) (*pb.Void, error) {
	mydao, err := c.daoCache.ForTenant(in.Package.Tenant, tblPackage, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	err = mydao.Delete(in.Package.ID)
	if err != nil {
		return nil, err
	}
	c.broadcastDelete(ctx, in.Package)
	return &pb.Void{}, nil
}

func (c *PackageController) UpdateToPackagesStr(_ *pb.UpdateToPackagesStrRequest, r pb.Control_UpdateToPackagesStrServer) error {
	return c.bJobPackage.RcvAndDispatchUpdates(c.ctx, r)
}

func (c *PackageController) getPackages(tenant string) ([]*pb.JobPackage, error) {
	mydao, err := c.daoCache.ForTenant(tenant, tblPackage, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	ms, err := mydao.All()
	if err != nil {
		return nil, err
	}
	packages := protoutil.Slices[*pb.JobPackage](ms)
	return packages, nil
}

func (c *PackageController) getPackage(tenant string, id string) (*pb.JobPackage, error) {
	mydao, err := c.daoCache.ForTenant(tenant, tblPackage, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	ms, err := mydao.Get(id)
	if err != nil {
		return nil, err
	}
	if ms != nil {
		return (*ms).(*pb.JobPackage), nil
	}
	return nil, nil
}

func (c *PackageController) broadcastAdd(ctx context.Context, m *pb.JobPackage) {
	c.broadcast(ctx, m, pb.UpdateType_New)
}

func (c *PackageController) broadcastUpdate(ctx context.Context, m *pb.JobPackage) {
	c.broadcast(ctx, m, pb.UpdateType_Update)
}

func (c *PackageController) broadcastDelete(ctx context.Context, m *pb.JobPackage) {
	c.broadcast(ctx, m, pb.UpdateType_Delete)
}

func (c *PackageController) broadcast(ctx context.Context, m *pb.JobPackage, utype pb.UpdateType) {
	c.bJobPackage.Broadcast(ctx, m, utype)
}
