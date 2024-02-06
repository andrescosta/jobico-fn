package client

import (
	"context"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/internal/api/types"
	rpc "google.golang.org/grpc"
)

type Queue struct {
	addr string
	conn *rpc.ClientConn
	cli  pb.QueueClient
}

func NewQueue(ctx context.Context, d service.GrpcDialer) (*Queue, error) {
	addr := env.String("queue.host")
	conn, err := d.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}
	cli := pb.NewQueueClient(conn)
	return &Queue{
		addr: addr,
		conn: conn,
		cli:  cli,
	}, nil
}

func (c *Queue) Close() error {
	return c.conn.Close()
}

func (c *Queue) Dequeue(ctx context.Context, tenant string, queue string) ([]*pb.QueueItem, error) {
	request := pb.DequeueRequest{
		Queue:  queue,
		Tenant: tenant,
	}
	r, err := c.cli.Dequeue(ctx, &request)
	if err != nil {
		return nil, err
	}
	return r.Items, nil
}

func (c *Queue) Queue(ctx context.Context, queueRequest *pb.QueueRequest) error {
	if _, err := c.cli.Queue(ctx, queueRequest); err != nil {
		return err
	}
	return nil
}
