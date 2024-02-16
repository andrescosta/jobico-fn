package server

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/database"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/ctl/controller"
)

type Server struct {
	pb.UnimplementedControlServer
	db              *database.Database
	pkgControler    *controller.PackageController
	envControler    *controller.EnvironmentController
	tenantControler *controller.TenantController
	ctx             context.Context
}

func New(ctx context.Context, dbDir string, dbo database.Option) (*Server, error) {
	db, err := database.Open(ctx, dbDir, dbo)
	if err != nil {
		return nil, err
	}
	return &Server{
		db:              db,
		tenantControler: controller.NewTenantController(db),
		pkgControler:    controller.NewPackageController(ctx, db),
		envControler:    controller.NewEnvironmentController(ctx, db),
		ctx:             ctx,
	}, nil
}

func (c *Server) Close() error {
	err := errors.Join(c.tenantControler.Close())
	err = errors.Join(err, c.pkgControler.Close())
	err = errors.Join(err, c.envControler.Close())
	err = errors.Join(err, c.db.Close(c.ctx))
	return err
}

func (c *Server) Packages(_ context.Context, in *pb.PackagesRequest) (*pb.PackagesReply, error) {
	return c.pkgControler.GetPackages(in)
}

func (c *Server) AllPackages(_ context.Context, _ *pb.Void) (*pb.AllPackagesReply, error) {
	return c.pkgControler.GetAllPackages()
}

func (c *Server) AddPackage(ctx context.Context, in *pb.AddPackageRequest) (*pb.AddPackageReply, error) {
	return c.pkgControler.AddPackage(ctx, in)
}

func (c *Server) UpdatePackage(ctx context.Context, in *pb.UpdatePackageRequest) (*pb.Void, error) {
	return c.pkgControler.UpdatePackage(ctx, in)
}

func (c *Server) DeletePackage(ctx context.Context, in *pb.DeletePackageRequest) (*pb.Void, error) {
	return c.pkgControler.DeletePackage(ctx, in)
}

func (c *Server) Tenants(_ context.Context, in *pb.TenantsRequest) (*pb.TenantsReply, error) {
	return c.tenantControler.Tenants(in)
}

func (c *Server) AddTenant(_ context.Context, in *pb.AddTenantRequest) (*pb.AddTenantReply, error) {
	return c.tenantControler.AddTenant(in)
}

func (c *Server) AddEnvironment(_ context.Context, in *pb.AddEnvironmentRequest) (*pb.AddEnvironmentReply, error) {
	return c.envControler.AddEnvironment(in)
}

func (c *Server) UpdateEnvironment(_ context.Context, in *pb.UpdateEnvironmentRequest) (*pb.Void, error) {
	return c.envControler.UpdateEnvironment(in)
}

func (c *Server) Environment(_ context.Context, _ *pb.Void) (*pb.EnvironmentReply, error) {
	return c.envControler.GetEnvironment()
}

func (c *Server) UpdateToPackagesStr(req *pb.UpdateToPackagesStrRequest, ctl pb.Control_UpdateToPackagesStrServer) error {
	return c.pkgControler.UpdateToPackagesStr(req, ctl)
}

func (c *Server) UpdateToEnvironmentStr(req *pb.Void, ctl pb.Control_UpdateToEnvironmentStrServer) error {
	return c.envControler.UpdateToEnvironmentStr(req, ctl)
}
