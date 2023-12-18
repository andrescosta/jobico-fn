package controller

import (
	"context"

	"github.com/andrescosta/goico/pkg/convert"
	"github.com/andrescosta/goico/pkg/database"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/ctl/internal/dao"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/protobuf/proto"
)

const (
	tblPackage = "package"
)

type PackageController struct {
	daoCache         *dao.Cache
	bJobPackage      *grpchelper.GrpcBroadcaster[*pb.UpdateToPackagesStrReply, proto.Message]
	tenantController *TenantController
}

func NewPackageController(ctx context.Context, db *database.Database) *PackageController {
	return &PackageController{
		daoCache:         dao.NewCache(db),
		bJobPackage:      grpchelper.StartBroadcaster[*pb.UpdateToPackagesStrReply, proto.Message](ctx),
		tenantController: NewTenantController(db),
	}
}
func (c *PackageController) Close() {
	c.bJobPackage.Stop()
}

func (c *PackageController) GetPackages(ctx context.Context, in *pb.GetJobPackagesRequest) (*pb.GetJobPackagesReply, error) {
	if in.ID != nil {
		p, err := c.getPackage(ctx, in.Tenant, *in.ID)
		if err != nil {
			return nil, err
		}
		if p != nil {
			return &pb.GetJobPackagesReply{Packages: []*pb.JobPackage{p}}, nil
		}
		return &pb.GetJobPackagesReply{}, nil
	}
	packages, err := c.getPackages(ctx, in.Tenant)
	if err != nil {
		return nil, err
	}
	return &pb.GetJobPackagesReply{Packages: packages}, nil
}

func (c *PackageController) GetAllPackages(ctx context.Context, _ *pb.GetAllJobPackagesRequest) (*pb.GetAllJobPackagesReply, error) {
	ms, err := c.tenantController.getTenants(ctx)
	if err != nil {
		return nil, err
	}
	packages := make([]*pb.JobPackage, 0)
	for _, me := range ms {
		mydao, err := c.daoCache.GetForTenant(ctx, me.ID, tblPackage, &pb.JobPackage{})
		if err != nil {
			return nil, err
		}
		ms, err := mydao.All(ctx)
		if err != nil {
			return nil, err
		}
		ps := convert.Slices[proto.Message, *pb.JobPackage](ms)
		packages = append(packages, ps...)
	}
	return &pb.GetAllJobPackagesReply{Packages: packages}, nil
}

func (c *PackageController) AddPackage(ctx context.Context, in *pb.AddJobPackageRequest) (*pb.AddJobPackageReply, error) {
	mydao, err := c.daoCache.GetForTenant(ctx, in.Package.Tenant, tblPackage, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Package
	_, err = mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	c.broadcastAdd(ctx, in.Package)
	return &pb.AddJobPackageReply{Package: in.Package}, nil
}
func (c *PackageController) UpdatePackage(ctx context.Context, in *pb.UpdateJobPackageRequest) (*pb.UpdateJobPackageReply, error) {
	mydao, err := c.daoCache.GetForTenant(ctx, in.Package.Tenant, tblPackage, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Package
	err = mydao.Update(ctx, m)
	if err != nil {
		return nil, err
	}
	c.broadcastUpdate(ctx, in.Package)
	return &pb.UpdateJobPackageReply{}, nil
}

func (c *PackageController) DeletePackage(ctx context.Context, in *pb.DeleteJobPackageRequest) (*pb.DeleteJobPackageReply, error) {
	mydao, err := c.daoCache.GetForTenant(ctx, in.Package.Tenant, tblPackage, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	err = mydao.Delete(ctx, in.Package.ID)
	if err != nil {
		return nil, err
	}
	c.broadcastDelete(ctx, in.Package)
	return &pb.DeleteJobPackageReply{}, nil
}
func (c *PackageController) UpdateToPackagesStr(_ *pb.UpdateToPackagesStrRequest, r pb.Control_UpdateToPackagesStrServer) error {
	return c.bJobPackage.RcvAndDispatchUpdates(r)
}
func (c *PackageController) getPackages(ctx context.Context, tenant string) ([]*pb.JobPackage, error) {
	mydao, err := c.daoCache.GetForTenant(ctx, tenant, tblPackage, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	ms, err := mydao.All(ctx)
	if err != nil {
		return nil, err
	}
	packages := convert.Slices[proto.Message, *pb.JobPackage](ms)
	return packages, nil
}
func (c *PackageController) getPackage(ctx context.Context, tenant string, id string) (*pb.JobPackage, error) {
	mydao, err := c.daoCache.GetForTenant(ctx, tenant, tblPackage, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	ms, err := mydao.Get(ctx, id)
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
