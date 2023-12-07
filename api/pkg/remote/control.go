package remote

import (
	"context"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/grpc"
)

type ControlClient struct {
	serverAddr string
	conn       *grpc.ClientConn
	client     pb.ControlClient
}

func NewControlClient(ctx context.Context) (*ControlClient, error) {
	host := env.GetAsString("ctl.host")
	conn, err := service.Dial(ctx, host)
	if err != nil {
		return nil, err
	}
	client := pb.NewControlClient(conn)
	return &ControlClient{
		serverAddr: host,
		conn:       conn,
		client:     client,
	}, nil
}

func (c *ControlClient) Close() {
	c.conn.Close()
}

func (c *ControlClient) GetEnviroment(ctx context.Context) (*pb.Environment, error) {
	r, err := c.client.GetEnviroment(ctx, &pb.GetEnviromentRequest{})
	if err != nil {
		return nil, err
	}
	return r.Environment, nil

}

func (c *ControlClient) AddEnvironment(ctx context.Context, environment *pb.Environment) (*pb.Environment, error) {
	r, err := c.client.AddEnviroment(ctx, &pb.AddEnviromentRequest{Environment: environment})
	if err != nil {
		return nil, err
	}
	return r.Environment, nil
}
func (c *ControlClient) UpdateEnvironment(ctx context.Context, environment *pb.Environment) error {
	if _, err := c.client.UpdateEnviroment(ctx, &pb.UpdateEnviromentRequest{Environment: environment}); err != nil {
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
		ID:       id,
		TenantId: tenant,
	})
	if err != nil {
		return nil, err
	}
	return r.Packages, nil
}

func (c *ControlClient) GetAllPackages(ctx context.Context) ([]*pb.JobPackage, error) {
	r, err := c.client.GetAllPackages(ctx, &pb.GetAllJobPackagesRequest{})
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

func (c *ControlClient) UpdateToPackagesStr(ctx context.Context, resChan chan<- *pb.UpdateToPackagesStrReply) error {
	s, err := c.client.UpdateToPackagesStr(ctx, &pb.UpdateToPackagesStrRequest{})
	if err != nil {
		return err
	}
	return grpchelper.Recv(ctx, s, resChan)
}

func (c *ControlClient) UpdateToEnviromentStr(ctx context.Context, resChan chan<- *pb.UpdateToEnviromentStrReply) error {
	s, err := c.client.UpdateToEnviromentStr(ctx, &pb.UpdateToEnviromentStrRequest{})
	if err != nil {
		return err
	}
	return grpchelper.Recv(ctx, s, resChan)
}
