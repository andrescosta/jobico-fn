package client

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/broadcaster"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	rpc "google.golang.org/grpc"
)

type Ctl struct {
	addr         string
	conn         *rpc.ClientConn
	cli          pb.ControlClient
	bcJobPackage *broadcaster.Broadcaster[*pb.UpdateToPackagesStrReply]
	bcEnvUpdates *broadcaster.Broadcaster[*pb.UpdateToEnvironmentStrReply]
}

var ErrCtlHostAddr = errors.New("the control service address was not specified in the env file using ctl.host")

func NewCtl(ctx context.Context, d service.GrpcDialer) (*Ctl, error) {
	host := env.StringOrNil("ctl.host")
	if host == nil {
		return nil, ErrCtlHostAddr
	}
	conn, err := d.Dial(ctx, *host)
	if err != nil {
		return nil, err
	}
	cli := pb.NewControlClient(conn)
	return &Ctl{
		addr: *host,
		conn: conn,
		cli:  cli,
	}, nil
}

func (c *Ctl) Close() error {
	errs := make([]error, 0)
	if c.bcEnvUpdates != nil {
		err := c.bcEnvUpdates.Stop()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if c.bcJobPackage != nil {
		err := c.bcJobPackage.Stop()
		if err != nil {
			errs = append(errs, err)
		}
	}
	err := c.conn.Close()
	if err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (c *Ctl) Environment(ctx context.Context) (*pb.Environment, error) {
	r, err := c.cli.Environment(ctx, &pb.Void{})
	if err != nil {
		return nil, err
	}
	return r.Environment, nil
}

func (c *Ctl) AddEnvironment(ctx context.Context, environment *pb.Environment) (*pb.Environment, error) {
	r, err := c.cli.AddEnvironment(ctx, &pb.AddEnvironmentRequest{Environment: environment})
	if err != nil {
		return nil, err
	}
	return r.Environment, nil
}

func (c *Ctl) UpdateEnvironment(ctx context.Context, environment *pb.Environment) error {
	if _, err := c.cli.UpdateEnvironment(ctx, &pb.UpdateEnvironmentRequest{Environment: environment}); err != nil {
		return err
	}
	return nil
}

func (c *Ctl) Tenants(ctx context.Context) ([]*pb.Tenant, error) {
	return c.tenants(ctx, nil)
}

func (c *Ctl) Tenant(ctx context.Context, id *string) ([]*pb.Tenant, error) {
	return c.tenants(ctx, id)
}

func (c *Ctl) tenants(ctx context.Context, id *string) ([]*pb.Tenant, error) {
	r, err := c.cli.Tenants(ctx, &pb.TenantsRequest{ID: id})
	if err != nil {
		return nil, err
	}
	return r.Tenants, nil
}

func (c *Ctl) AddTenant(ctx context.Context, tenant *pb.Tenant) (*pb.Tenant, error) {
	r, err := c.cli.AddTenant(ctx, &pb.AddTenantRequest{Tenant: tenant})
	if err != nil {
		return nil, err
	}
	return r.Tenant, nil
}

func (c *Ctl) Packages(ctx context.Context, tenant string) ([]*pb.JobPackage, error) {
	return c.packages(ctx, tenant, nil)
}

func (c *Ctl) Package(ctx context.Context, tenant string, id *string) ([]*pb.JobPackage, error) {
	return c.packages(ctx, tenant, id)
}

func (c *Ctl) packages(ctx context.Context, tenant string, id *string) ([]*pb.JobPackage, error) {
	r, err := c.cli.Packages(ctx, &pb.PackagesRequest{
		ID:     id,
		Tenant: tenant,
	})
	if err != nil {
		return nil, err
	}
	return r.Packages, nil
}

func (c *Ctl) AllPackages(ctx context.Context) ([]*pb.JobPackage, error) {
	r, err := c.cli.AllPackages(ctx, &pb.Void{})
	if err != nil {
		return nil, err
	}
	return r.Packages, nil
}

func (c *Ctl) AddPackage(ctx context.Context, pkg *pb.JobPackage) (*pb.JobPackage, error) {
	r, err := c.cli.AddPackage(ctx, &pb.AddPackageRequest{Package: pkg})
	if err != nil {
		return nil, err
	}
	return r.Package, nil
}

func (c *Ctl) UpdatePackage(ctx context.Context, pkg *pb.JobPackage) error {
	_, err := c.cli.UpdatePackage(ctx, &pb.UpdatePackageRequest{Package: pkg})
	if err != nil {
		return err
	}
	return nil
}

func (c *Ctl) DeletePackage(ctx context.Context, pkg *pb.JobPackage) error {
	_, err := c.cli.DeletePackage(ctx, &pb.DeletePackageRequest{Package: pkg})
	if err != nil {
		return err
	}
	return nil
}

func (c *Ctl) ListenerForEnvironmentUpdates(ctx context.Context) (*broadcaster.Listener[*pb.UpdateToEnvironmentStrReply], error) {
	if c.bcEnvUpdates == nil {
		if err := c.startListenEnvironmentUpdates(ctx); err != nil {
			return nil, err
		}
	}
	return c.bcEnvUpdates.Subscribe()
}

func (c *Ctl) startListenEnvironmentUpdates(ctx context.Context) error {
	cb := broadcaster.Start[*pb.UpdateToEnvironmentStrReply](ctx)
	c.bcEnvUpdates = cb
	s, err := c.cli.UpdateToEnvironmentStr(ctx, &pb.Void{})
	if err != nil {
		return err
	}
	go func() {
		_ = grpchelper.Listen(ctx, s, cb)
	}()
	return nil
}

func (c *Ctl) ListenerForPackageUpdates(ctx context.Context) (*broadcaster.Listener[*pb.UpdateToPackagesStrReply], error) {
	if c.bcJobPackage == nil {
		if err := c.startListenerForPackageUpdates(ctx); err != nil {
			return nil, err
		}
	}
	return c.bcJobPackage.Subscribe()
}

func (c *Ctl) startListenerForPackageUpdates(ctx context.Context) error {
	cb := broadcaster.Start[*pb.UpdateToPackagesStrReply](ctx)
	c.bcJobPackage = cb
	s, err := c.cli.UpdateToPackagesStr(ctx, &pb.UpdateToPackagesStrRequest{})
	if err != nil {
		return err
	}
	go func() {
		_ = grpchelper.Listen(ctx, s, cb)
	}()
	return nil
}
