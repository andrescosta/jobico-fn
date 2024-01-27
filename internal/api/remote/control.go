package remote

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

type CtlClient struct {
	serverAddr            string
	conn                  *rpc.ClientConn
	client                pb.ControlClient
	broadcasterJobPackage *broadcaster.Broadcaster[*pb.UpdateToPackagesStrReply]
	broadcasterEnvUpdates *broadcaster.Broadcaster[*pb.UpdateToEnvironmentStrReply]
}

var ErrCtlHostAddr = errors.New("the control service address was not specified in the env file using ctl.host")

func NewCtlClient(ctx context.Context, d service.GrpcDialer) (*CtlClient, error) {
	host := env.StringOrNil("ctl.host")
	if host == nil {
		return nil, ErrCtlHostAddr
	}
	conn, err := d.Dial(ctx, *host)
	if err != nil {
		return nil, err
	}
	client := pb.NewControlClient(conn)
	return &CtlClient{
		serverAddr: *host,
		conn:       conn,
		client:     client,
	}, nil
}

func (c *CtlClient) Close() error {
	if c.broadcasterEnvUpdates != nil {
		c.broadcasterEnvUpdates.Stop()
	}
	if c.broadcasterJobPackage != nil {
		c.broadcasterJobPackage.Stop()
	}
	return c.conn.Close()
}

func (c *CtlClient) GetEnvironment(ctx context.Context) (*pb.Environment, error) {
	r, err := c.client.GetEnvironment(ctx, &pb.Void{})
	if err != nil {
		return nil, err
	}
	return r.Environment, nil
}

func (c *CtlClient) AddEnvironment(ctx context.Context, environment *pb.Environment) (*pb.Environment, error) {
	r, err := c.client.AddEnvironment(ctx, &pb.AddEnvironmentRequest{Environment: environment})
	if err != nil {
		return nil, err
	}
	return r.Environment, nil
}

func (c *CtlClient) UpdateEnvironment(ctx context.Context, environment *pb.Environment) error {
	if _, err := c.client.UpdateEnvironment(ctx, &pb.UpdateEnvironmentRequest{Environment: environment}); err != nil {
		return err
	}
	return nil
}

func (c *CtlClient) GetTenants(ctx context.Context) ([]*pb.Tenant, error) {
	return c.getTenants(ctx, nil)
}

func (c *CtlClient) GetTenant(ctx context.Context, id *string) ([]*pb.Tenant, error) {
	return c.getTenants(ctx, id)
}

func (c *CtlClient) getTenants(ctx context.Context, id *string) ([]*pb.Tenant, error) {
	r, err := c.client.GetTenants(ctx, &pb.GetTenantsRequest{ID: id})
	if err != nil {
		return nil, err
	}
	return r.Tenants, nil
}

func (c *CtlClient) AddTenant(ctx context.Context, tenant *pb.Tenant) (*pb.Tenant, error) {
	r, err := c.client.AddTenant(ctx, &pb.AddTenantRequest{Tenant: tenant})
	if err != nil {
		return nil, err
	}
	return r.Tenant, nil
}

func (c *CtlClient) GetPackages(ctx context.Context, tenant string) ([]*pb.JobPackage, error) {
	return c.getPackages(ctx, tenant, nil)
}

func (c *CtlClient) GetPackage(ctx context.Context, tenant string, id *string) ([]*pb.JobPackage, error) {
	return c.getPackages(ctx, tenant, id)
}

func (c *CtlClient) getPackages(ctx context.Context, tenant string, id *string) ([]*pb.JobPackage, error) {
	r, err := c.client.GetPackages(ctx, &pb.GetJobPackagesRequest{
		ID:     id,
		Tenant: tenant,
	})
	if err != nil {
		return nil, err
	}
	return r.Packages, nil
}

func (c *CtlClient) GetAllPackages(ctx context.Context) ([]*pb.JobPackage, error) {
	r, err := c.client.GetAllPackages(ctx, &pb.Void{})
	if err != nil {
		return nil, err
	}
	return r.Packages, nil
}

func (c *CtlClient) AddPackage(ctx context.Context, package1 *pb.JobPackage) (*pb.JobPackage, error) {
	r, err := c.client.AddPackage(ctx, &pb.AddJobPackageRequest{Package: package1})
	if err != nil {
		return nil, err
	}
	return r.Package, nil
}

func (c *CtlClient) UpdatePackage(ctx context.Context, package1 *pb.JobPackage) error {
	_, err := c.client.UpdatePackage(ctx, &pb.UpdateJobPackageRequest{Package: package1})
	if err != nil {
		return err
	}
	return nil
}

func (c *CtlClient) DeletePackage(ctx context.Context, package1 *pb.JobPackage) error {
	_, err := c.client.DeletePackage(ctx, &pb.DeleteJobPackageRequest{Package: package1})
	if err != nil {
		return err
	}
	return nil
}

func (c *CtlClient) ListenerForEnvironmentUpdates(ctx context.Context) (*broadcaster.Listener[*pb.UpdateToEnvironmentStrReply], error) {
	if c.broadcasterEnvUpdates == nil {
		if err := c.startListenEnvironmentUpdates(ctx); err != nil {
			return nil, err
		}
	}
	return c.broadcasterEnvUpdates.Subscribe()
}

func (c *CtlClient) startListenEnvironmentUpdates(ctx context.Context) error {
	cb := broadcaster.Start[*pb.UpdateToEnvironmentStrReply](ctx)
	c.broadcasterEnvUpdates = cb
	s, err := c.client.UpdateToEnvironmentStr(ctx, &pb.Void{})
	if err != nil {
		return err
	}
	go func() {
		_ = grpchelper.Listen(ctx, s, cb)
	}()
	return nil
}

func (c *CtlClient) ListenerForPackageUpdates(ctx context.Context) (*broadcaster.Listener[*pb.UpdateToPackagesStrReply], error) {
	if c.broadcasterJobPackage == nil {
		if err := c.startListenerForPackageUpdates(ctx); err != nil {
			return nil, err
		}
	}
	return c.broadcasterJobPackage.Subscribe()
}

func (c *CtlClient) startListenerForPackageUpdates(ctx context.Context) error {
	cb := broadcaster.Start[*pb.UpdateToPackagesStrReply](ctx)
	c.broadcasterJobPackage = cb
	s, err := c.client.UpdateToPackagesStr(ctx, &pb.UpdateToPackagesStrRequest{})
	if err != nil {
		return err
	}
	go func() {
		_ = grpchelper.Listen(ctx, s, cb)
	}()
	return nil
}
