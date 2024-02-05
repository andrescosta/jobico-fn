package server

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/database"
	"github.com/andrescosta/goico/pkg/env"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/ctl/controller"
)

type Server struct {
	pb.UnimplementedControlServer
	db         *database.Database
	pkgCont    *controller.PackageController
	envCont    *controller.EnvironmentController
	tenantCont *controller.TenantController
}

func New(ctx context.Context, dbFileName string, dbo database.Option) (*Server, error) {
	dbPath := env.WorkdirPlus(dbFileName)
	db, err := database.Open(dbPath, dbo)
	if err != nil {
		return nil, err
	}
	return &Server{
		db:         db,
		tenantCont: controller.NewTenantController(db),
		pkgCont:    controller.NewPackageController(ctx, db),
		envCont:    controller.NewEnvironmentController(ctx, db),
	}, nil
}

func (c *Server) Close() error {
	var err error
	err = errors.Join(err, c.tenantCont.Close())
	err = errors.Join(err, c.pkgCont.Close())
	err = errors.Join(err, c.envCont.Close())
	err = errors.Join(err, c.db.Close())
	return err
}

func (c *Server) Packages(_ context.Context, in *pb.PackagesRequest) (*pb.PackagesReply, error) {
	return c.pkgCont.GetPackages(in)
}

func (c *Server) AllPackages(_ context.Context, _ *pb.Void) (*pb.AllPackagesReply, error) {
	return c.pkgCont.GetAllPackages()
}

func (c *Server) AddPackage(ctx context.Context, in *pb.AddPackageRequest) (*pb.AddPackageReply, error) {
	return c.pkgCont.AddPackage(ctx, in)
}

func (c *Server) UpdatePackage(ctx context.Context, in *pb.UpdatePackageRequest) (*pb.Void, error) {
	return c.pkgCont.UpdatePackage(ctx, in)
}

func (c *Server) DeletePackage(ctx context.Context, in *pb.DeletePackageRequest) (*pb.Void, error) {
	return c.pkgCont.DeletePackage(ctx, in)
}

func (c *Server) Tenants(_ context.Context, in *pb.TenantsRequest) (*pb.TenantsReply, error) {
	return c.tenantCont.Tenants(in)
}

func (c *Server) AddTenant(_ context.Context, in *pb.AddTenantRequest) (*pb.AddTenantReply, error) {
	return c.tenantCont.AddTenant(in)
}

func (c *Server) AddEnvironment(_ context.Context, in *pb.AddEnvironmentRequest) (*pb.AddEnvironmentReply, error) {
	return c.envCont.AddEnvironment(in)
}

func (c *Server) UpdateEnvironment(_ context.Context, in *pb.UpdateEnvironmentRequest) (*pb.Void, error) {
	return c.envCont.UpdateEnvironment(in)
}

func (c *Server) Environment(_ context.Context, _ *pb.Void) (*pb.EnvironmentReply, error) {
	return c.envCont.GetEnvironment()
}

func (c *Server) UpdateToPackagesStr(req *pb.UpdateToPackagesStrRequest, ctl pb.Control_UpdateToPackagesStrServer) error {
	return c.pkgCont.UpdateToPackagesStr(req, ctl)
}

func (c *Server) UpdateToEnvironmentStr(req *pb.Void, ctl pb.Control_UpdateToEnvironmentStrServer) error {
	return c.envCont.UpdateToEnvironmentStr(req, ctl)
}
