package server

import (
	"context"

	"github.com/andrescosta/goico/pkg/database"
	"github.com/andrescosta/goico/pkg/env"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/ctl/controller"
)

type Server struct {
	pb.UnimplementedControlServer
	db         *database.Database
	pkgCont    *controller.PackageController
	envCont    *controller.EnvironmentController
	tenantCont *controller.TenantController
}

func New(ctx context.Context, dbFileName string) (*Server, error) {
	dbPath := env.ElemInWorkDir(dbFileName)
	db, err := database.Open(dbPath)
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
	c.tenantCont.Close()
	c.pkgCont.Close()
	c.envCont.Close()
	return c.db.Close()
}

func (c *Server) GetPackages(ctx context.Context, in *pb.GetJobPackagesRequest) (*pb.GetJobPackagesReply, error) {
	return c.pkgCont.GetPackages(ctx, in)
}

func (c *Server) GetAllPackages(ctx context.Context, in *pb.GetAllJobPackagesRequest) (*pb.GetAllJobPackagesReply, error) {
	return c.pkgCont.GetAllPackages(ctx, in)
}

func (c *Server) AddPackage(ctx context.Context, in *pb.AddJobPackageRequest) (*pb.AddJobPackageReply, error) {
	return c.pkgCont.AddPackage(ctx, in)
}

func (c *Server) UpdatePackage(ctx context.Context, in *pb.UpdateJobPackageRequest) (*pb.UpdateJobPackageReply, error) {
	return c.pkgCont.UpdatePackage(ctx, in)
}

func (c *Server) DeletePackage(ctx context.Context, in *pb.DeleteJobPackageRequest) (*pb.DeleteJobPackageReply, error) {
	return c.pkgCont.DeletePackage(ctx, in)
}

func (c *Server) GetTenants(ctx context.Context, in *pb.GetTenantsRequest) (*pb.GetTenantsReply, error) {
	return c.tenantCont.GetTenants(ctx, in)
}

func (c *Server) AddTenant(ctx context.Context, in *pb.AddTenantRequest) (*pb.AddTenantReply, error) {
	return c.tenantCont.AddTenant(ctx, in)
}

func (c *Server) AddEnvironment(ctx context.Context, in *pb.AddEnvironmentRequest) (*pb.AddEnvironmentReply, error) {
	return c.envCont.AddEnvironment(ctx, in)
}

func (c *Server) UpdateEnvironment(ctx context.Context, in *pb.UpdateEnvironmentRequest) (*pb.UpdateEnvironmentReply, error) {
	return c.envCont.UpdateEnvironment(ctx, in)
}

func (c *Server) GetEnvironment(ctx context.Context, in *pb.GetEnvironmentRequest) (*pb.GetEnvironmentReply, error) {
	return c.envCont.GetEnvironment(ctx, in)
}

func (c *Server) UpdateToPackagesStr(req *pb.UpdateToPackagesStrRequest, ctl pb.Control_UpdateToPackagesStrServer) error {
	return c.pkgCont.UpdateToPackagesStr(req, ctl)
}

func (c *Server) UpdateToEnvironmentStr(req *pb.UpdateToEnvironmentStrRequest, ctl pb.Control_UpdateToEnvironmentStrServer) error {
	return c.envCont.UpdateToEnvironmentStr(req, ctl)
}
