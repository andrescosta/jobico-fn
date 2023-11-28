package remote

import (
	"context"

	"github.com/andrescosta/goico/pkg/env"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RecorderClient struct {
	serverAddr string
}

func NewRecorderClient() *RecorderClient {
	return &RecorderClient{
		serverAddr: env.GetAsString("recorder.host"),
	}
}

func (c *RecorderClient) dial() (*grpc.ClientConn, error) {
	ops := grpc.WithTransportCredentials(insecure.NewCredentials())
	return grpc.Dial(c.serverAddr, ops)

}

func (c *RecorderClient) GetJobExecutions(ctx context.Context, tenant string, lines int32, resChan chan<- string) error {
	logger := zerolog.Ctx(ctx)

	conn, err := c.dial()
	if err != nil {
		return err
	}
	defer conn.Close()
	recorderClient := pb.NewRecorderClient(conn)

	rj, err := recorderClient.GetJobExecutions(ctx, &pb.GetJobExecutionsRequest{
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
				logger.Err(err).Msg("error getting message")
			} else {
				for _, r := range ress.Result {
					resChan <- r
				}
			}
		}
	}
}

func (c *RecorderClient) AddJobExecution(ctx context.Context, ex *pb.JobExecution) error {
	conn, err := c.dial()
	if err != nil {
		return err
	}
	defer conn.Close()
	recorderClient := pb.NewRecorderClient(conn)
	_, err = recorderClient.AddJobExecution(ctx, &pb.AddJobExecutionRequest{
		Execution: ex,
	})
	if err != nil {
		return err
	}
	return nil
}
