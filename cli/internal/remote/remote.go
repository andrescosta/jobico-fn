package remote

import (
	"context"

	pb "github.com/andrescosta/workflew/api/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ControlClient struct {
	serverAddr string
}

func NewControlClient(serverAddr string) ControlClient {
	return ControlClient{
		serverAddr: serverAddr,
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

func (c *ControlClient) GetMerchants(ctx context.Context) ([]*pb.Merchant, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	r, err := client.GetMechants(ctx, &pb.GetMerchantsRequest{})
	if err != nil {
		return nil, err
	}
	return r.Merchants, nil
}

func (c *ControlClient) AddMerchant(ctx context.Context, merchant *pb.Merchant) (*pb.Merchant, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewControlClient(conn)
	r, err := client.AddMerchant(ctx, &pb.AddMerchantRequest{Merchant: merchant})
	if err != nil {
		return nil, err
	}
	return r.Merchant, nil
}

func (c *ControlClient) GetQueueDefs(ctx context.Context) ([]*pb.QueueDef, error) {
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
}
