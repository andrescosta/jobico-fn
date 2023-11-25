package remote

import (
	"context"

	"github.com/andrescosta/goico/pkg/env"
	pb "github.com/andrescosta/workflew/api/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QueueClient struct {
	serverAddr string
}

func NewQueueClient() *QueueClient {
	return &QueueClient{
		serverAddr: env.GetAsString("queue.host"),
	}
}

func (c *QueueClient) dial() (*grpc.ClientConn, error) {
	ops := grpc.WithTransportCredentials(insecure.NewCredentials())
	return grpc.Dial(c.serverAddr, ops)

}

func (c *QueueClient) Dequeue(ctx context.Context, tenant string, queue string) ([]*pb.QueueItem, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewQueueClient(conn)
	request := pb.DequeueRequest{
		QueueId:  queue,
		TenantId: tenant,
	}
	r, err := client.Dequeue(ctx, &request)
	if err != nil {
		return nil, err
	}
	return r.Items, nil
}
