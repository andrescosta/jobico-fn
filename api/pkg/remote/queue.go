package remote

import (
	"context"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/api/types"
	"google.golang.org/grpc"
)

type QueueClient struct {
	serverAddr string
	conn       *grpc.ClientConn
	client     pb.QueueClient
}

func NewQueueClient() (*QueueClient, error) {
	host := env.GetAsString("queue.host")
	conn, err := service.Dial(host)
	if err != nil {
		return nil, err
	}
	client := pb.NewQueueClient(conn)
	return &QueueClient{
		serverAddr: host,
		conn:       conn,
		client:     client,
	}, nil
}

func (c *QueueClient) Close() {
	c.conn.Close()
}

func (c *QueueClient) Dequeue(ctx context.Context, tenant string, queue string) ([]*pb.QueueItem, error) {
	request := pb.DequeueRequest{
		QueueId:  queue,
		TenantId: tenant,
	}
	r, err := c.client.Dequeue(ctx, &request)
	if err != nil {
		return nil, err
	}
	return r.Items, nil
}

func (c *QueueClient) Queue(ctx context.Context, queueRequest *pb.QueueRequest) error {
	if _, err := c.client.Queue(ctx, queueRequest); err != nil {
		return err
	}
	return nil
}
