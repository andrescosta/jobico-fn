package remote

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/broadcaster"
	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	rpc "google.golang.org/grpc"
)

type ControlClient struct {
	serverAddr            string
	conn                  *rpc.ClientConn
	client                pb.ControlClient
	broadcasterJobPackage *broadcaster.Broadcaster[*pb.UpdateToPackagesStrReply]
	broadcasterEnvUpdates *broadcaster.Broadcaster[*pb.UpdateToEnvironmentStrReply]
}

var ErrCtlHostAddr = errors.New("the control service address was not specified in the env file using ctl.host")

func NewControlClient(ctx context.Context, d service.GrpcDialer) (*ControlClient, error) {
	host := env.StringOrNil("ctl.host")
	if host == nil {
		return nil, ErrCtlHostAddr
	}
	conn, err := d.Dial(ctx, *host)
	if err != nil {
		return nil, err
	}
	client := pb.NewControlClient(conn)
	return &ControlClient{
		serverAddr: *host,
		conn:       conn,
		client:     client,
	}, nil
}

func (c *ControlClient) Close() {
	_ = c.conn.Close()
}

func (c *ControlClient) GetEnvironment(ctx context.Context) (*pb.Environment, error) {
	r, err := c.client.GetEnvironment(ctx, &pb.Void{})
	if err != nil {
		return nil, err
	}
	return r.Environment, nil
}

func (c *ControlClient) AddEnvironment(ctx context.Context, environment *pb.Environment) (*pb.Environment, error) {
	r, err := c.client.AddEnvironment(ctx, &pb.AddEnvironmentRequest{Environment: environment})
	if err != nil {
		return nil, err
	}
	return r.Environment, nil
}

func (c *ControlClient) UpdateEnvironment(ctx context.Context, environment *pb.Environment) error {
	if _, err := c.client.UpdateEnvironment(ctx, &pb.UpdateEnvironmentRequest{Environment: environment}); err != nil {
		return err
	}
	return nil
}

func (c *ControlClient) GetTenants(ctx context.Context) ([]*pb.Tenant, error) {
	return c.getTenants(ctx, nil)
}

func (c *ControlClient) GetTenant(ctx context.Context, id *string) ([]*pb.Tenant, error) {
	return c.getTenants(ctx, id)
}

func (c *ControlClient) getTenants(ctx context.Context, id *string) ([]*pb.Tenant, error) {
	r, err := c.client.GetTenants(ctx, &pb.GetTenantsRequest{ID: id})
	if err != nil {
		return nil, err
	}
	return r.Tenants, nil
}

func (c *ControlClient) AddTenant(ctx context.Context, tenant *pb.Tenant) (*pb.Tenant, error) {
	r, err := c.client.AddTenant(ctx, &pb.AddTenantRequest{Tenant: tenant})
	if err != nil {
		return nil, err
	}
	return r.Tenant, nil
}

func (c *ControlClient) GetPackages(ctx context.Context, tenant string) ([]*pb.JobPackage, error) {
	return c.getPackages(ctx, tenant, nil)
}

func (c *ControlClient) GetPackage(ctx context.Context, tenant string, id *string) ([]*pb.JobPackage, error) {
	return c.getPackages(ctx, tenant, id)
}

func (c *ControlClient) getPackages(ctx context.Context, tenant string, id *string) ([]*pb.JobPackage, error) {
	r, err := c.client.GetPackages(ctx, &pb.GetJobPackagesRequest{
		ID:     id,
		Tenant: tenant,
	})
	if err != nil {
		return nil, err
	}
	return r.Packages, nil
}

func (c *ControlClient) GetAllPackages(ctx context.Context) ([]*pb.JobPackage, error) {
	r, err := c.client.GetAllPackages(ctx, &pb.Void{})
	if err != nil {
		return nil, err
	}
	return r.Packages, nil
}

func (c *ControlClient) AddPackage(ctx context.Context, package1 *pb.JobPackage) (*pb.JobPackage, error) {
	r, err := c.client.AddPackage(ctx, &pb.AddJobPackageRequest{Package: package1})
	if err != nil {
		return nil, err
	}
	return r.Package, nil
}

func (c *ControlClient) UpdatePackage(ctx context.Context, package1 *pb.JobPackage) error {
	_, err := c.client.UpdatePackage(ctx, &pb.UpdateJobPackageRequest{Package: package1})
	if err != nil {
		return err
	}
	return nil
}

func (c *ControlClient) DeletePackage(ctx context.Context, package1 *pb.JobPackage) error {
	_, err := c.client.DeletePackage(ctx, &pb.DeleteJobPackageRequest{Package: package1})
	if err != nil {
		return err
	}
	return nil
}

func (c *ControlClient) ListenerForEnvironmentUpdates(ctx context.Context) (*broadcaster.Listener[*pb.UpdateToEnvironmentStrReply], error) {
	if c.broadcasterEnvUpdates == nil {
		if err := c.startListenEnvironmentUpdates(ctx); err != nil {
			return nil, err
		}
	}
	return c.broadcasterEnvUpdates.Subscribe()
}

func (c *ControlClient) startListenEnvironmentUpdates(ctx context.Context) error {
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

func (c *ControlClient) ListenerForPackageUpdates(ctx context.Context) (*broadcaster.Listener[*pb.UpdateToPackagesStrReply], error) {
	if c.broadcasterJobPackage == nil {
		if err := c.startListenerForPackageUpdates(ctx); err != nil {
			return nil, err
		}
	}
	return c.broadcasterJobPackage.Subscribe()
}

func (c *ControlClient) startListenerForPackageUpdates(ctx context.Context) error {
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
