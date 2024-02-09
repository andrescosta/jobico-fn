package client

import (
	"context"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/rs/zerolog"
	rpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Recorder struct {
	addr string
	conn *rpc.ClientConn
	cli  pb.RecorderClient
}

func NewRecorder(ctx context.Context, dialer service.GrpcDialer) (*Recorder, error) {
	addr := env.String("recorder.host")
	conn, err := dialer.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}
	cli := pb.NewRecorderClient(conn)
	return &Recorder{
		addr: addr,
		conn: conn,
		cli:  cli,
	}, nil
}

func (c *Recorder) Close() error {
	return c.conn.Close()
}

func (c *Recorder) StreamJobExecutions(ctx context.Context, lines int32, resChan chan<- string) error {
	logger := zerolog.Ctx(ctx)
	rj, err := c.cli.GetJobExecutionsStr(ctx, &pb.JobExecutionsRequest{
		Lines: &lines,
	})
	if err != nil {
		return err
	}
	for {
		select {
		case <-rj.Context().Done():
			if err := rj.CloseSend(); err != nil {
				logger.Warn().AnErr("err", err).Msg("Recorder Client: error closing client stream")
			}
			return nil
		default:
			ress, err := rj.Recv()
			if err != nil {
				if status.Code(err) != codes.Canceled {
					logger.Warn().Msgf("error getting message %s", err)
				}
			} else {
				for _, r := range ress.Result {
					resChan <- r
				}
			}
		}
	}
}

func (c *Recorder) AddJobExecution(ctx context.Context, ex *pb.JobExecution) error {
	if _, err := c.cli.AddJobExecution(ctx, &pb.AddJobExecutionRequest{Execution: ex}); err != nil {
		return err
	}
	return nil
}

func (c *Recorder) JobExecutions(ctx context.Context, tenant string, lines int32) ([]string, error) {
	res, err := c.cli.JobExecutions(ctx, &pb.JobExecutionsRequest{
		Tenant: &tenant,
		Lines:  &lines,
	})
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}
