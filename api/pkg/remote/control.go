package remote

import (
	"context"

	"github.com/andrescosta/goico/pkg/env"
	pb "github.com/andrescosta/workflew/api/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ControlClient struct {
	serverAddr string
}

func NewControlClient() *ControlClient {
	return &ControlClient{
		serverAddr: env.GetAsString("ctl.host"),
	}
}

func (c *ControlClient) dial() (*grpc.ClientConn, error) {
	ops := grpc.WithTransportCredentials(insecure.NewCredentials())
	return grpc.Dial(c.serverAddr, ops)

}

func (c *ControlClient) GetEnviroment(ctx context.Context) (*pb.Environment, error) {

	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	r, err := client.GetEnviroment(ctx, &pb.GetEnviromentRequest{})
	if err != nil {
		return nil, err
	}
	return r.Environment, nil

}

func (c *ControlClient) AddEnvironment(ctx context.Context, environment *pb.Environment) (*pb.Environment, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	r, err := client.AddEnviroment(ctx, &pb.AddEnviromentRequest{Environment: environment})
	if err != nil {
		return nil, err
	}
	return r.Environment, nil
}
func (c *ControlClient) UpdateEnvironment(ctx context.Context, environment *pb.Environment) error {
	conn, err := c.dial()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	_, err = client.UpdateEnviroment(ctx, &pb.UpdateEnviromentRequest{Environment: environment})
	if err != nil {
		return err
	}
	return nil
}

func (c *ControlClient) GetTenants(ctx context.Context) ([]*pb.Tenant, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	r, err := client.GetTenants(ctx, &pb.GetTenantsRequest{})
	if err != nil {
		return nil, err
	}
	return r.Tenants, nil
}

func (c *ControlClient) AddTenant(ctx context.Context, tenant *pb.Tenant) (*pb.Tenant, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	r, err := client.AddTenant(ctx, &pb.AddTenantRequest{Tenant: tenant})
	if err != nil {
		return nil, err
	}
	return r.Tenant, nil
}

func (c *ControlClient) GetPackages(ctx context.Context, tenant string) ([]*pb.JobPackage, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	r, err := client.GetPackages(ctx, &pb.GetJobPackagesRequest{
		TenantId: tenant,
	})
	if err != nil {
		return nil, err
	}
	return r.Packages, nil
}

func (c *ControlClient) GetAllPackages(ctx context.Context) ([]*pb.JobPackage, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	r, err := client.GetAllPackages(ctx, &pb.GetAllJobPackagesRequest{})
	if err != nil {
		return nil, err
	}
	return r.Packages, nil
}

func (c *ControlClient) AddPackage(ctx context.Context, package1 *pb.JobPackage) (*pb.JobPackage, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	r, err := client.AddPackage(ctx, &pb.AddJobPackageRequest{Package: package1})
	if err != nil {
		return nil, err
	}
	return r.Package, nil
}

/*
	func (c *ControlClient) GetEventDefs(ctx context.Context) ([]*pb.EventDef, error) {
		conn, err := c.dial()
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := pb.NewControlClient(conn)
		r, err := client.GetEventDefs(ctx, &pb.GetEventDefsRequest{})
		if err != nil {
			return nil, err
		}
		return r.Events, nil
	}

	func (c *ControlClient) AddEventDef(ctx context.Context, event *pb.EventDef) (*pb.EventDef, error) {
		conn, err := c.dial()
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		client := pb.NewControlClient(conn)
		r, err := client.AddEventDef(ctx, &pb.AddEventDefRequest{Event: event})
		if err != nil {
			return nil, err
		}
		return r.Event, nil
	}
*/
/*func (c *ControlClient) GetQueueDefs(ctx context.Context) ([]*pb.QueueDef, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	r, err := client.GetQueueDefs(ctx, &pb.GetQueueDefsRequest{})
	if err != nil {
		return nil, err
	}
	return r.Queues, nil
}

func (c *ControlClient) AddQueueDef(ctx context.Context, queue *pb.QueueDef) (*pb.QueueDef, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	r, err := client.AddQueueDef(ctx, &pb.AddQueueDefRequest{Queue: queue})
	if err != nil {
		return nil, err
	}
	return r.Queue, nil
}*/
