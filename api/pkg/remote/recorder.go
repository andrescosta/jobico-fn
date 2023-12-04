package remote

import (
	"context"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type RecorderClient struct {
	serverAddr string
	conn       *grpc.ClientConn
	client     pb.RecorderClient
}

func NewRecorderClient() (*RecorderClient, error) {
	addr := env.GetAsString("recorder.host")
	conn, err := service.Dial(addr)
	if err != nil {
		return nil, err
	}
	client := pb.NewRecorderClient(conn)

	return &RecorderClient{
		serverAddr: addr,
		conn:       conn,
		client:     client,
	}, nil
}

func (c *RecorderClient) Close() {
	c.conn.Close()
}

func (c *RecorderClient) GetJobExecutions(ctx context.Context, tenant string, lines int32, resChan chan<- string) error {
	logger := zerolog.Ctx(ctx)
	rj, err := c.client.GetJobExecutions(ctx, &pb.GetJobExecutionsRequest{
		Lines: &lines,
	})
	if err != nil {
		return err
	}
	for {
		select {
		case <-rj.Context().Done():
			rj.CloseSend()
			return nil
		case <-ctx.Done():
			rj.CloseSend()
			return nil
		default:
			ress, err := rj.Recv()
			if err != nil {
				logger.Warn().Msgf("error getting message %s", err)
			} else {
				for _, r := range ress.Result {
					resChan <- r
				}
			}
		}
	}
}

func (c *RecorderClient) AddJobExecution(ctx context.Context, ex *pb.JobExecution) error {
	if _, err := c.client.AddJobExecution(ctx, &pb.AddJobExecutionRequest{Execution: ex}); err != nil {
		return err
	}
	return nil
}
