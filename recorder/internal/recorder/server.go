package recorder

import (
	"context"
	"io"

	iog "github.com/andrescosta/goico/pkg/io"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/nxadm/tail"
	"github.com/rs/zerolog"
)

type Server struct {
	pb.UnimplementedRecorderServer
	recorder *Recorder
	fullpath string
	ctx      context.Context
	//Repo *FileRepo

}

func NewServer(ctx context.Context, fullpath string) *Server {
	r := NewRecorder(fullpath)
	return &Server{
		recorder: r,
		fullpath: fullpath,
		ctx:      ctx,
	}
}

func (s *Server) AddJobExecution(ctx context.Context, r *pb.AddJobExecutionRequest) (*pb.AddJobExecutionReply, error) {
	s.recorder.AddExecution(r.Execution)
	return &pb.AddJobExecutionReply{}, nil
}

func (s *Server) GetJobExecutions(g *pb.GetJobExecutionsRequest, r pb.Recorder_GetJobExecutionsServer) error {
	logger := zerolog.Ctx(r.Context())
	seekInfo := &tail.SeekInfo{
		Offset: 0,
		Whence: io.SeekEnd,
	}
	if g.Lines != nil && *g.Lines > 0 {
		lines, err := iog.GetLastnLines(s.fullpath, int(*g.Lines), true, true)
		if err == nil {
			r.Send(&pb.GetJobExecutionsReply{
				Result: lines,
			})
		} else {
			logger.Warn().Msgf("error getting tail lines %s", err)
		}
	}
	tail, err := tail.TailFile(s.fullpath, tail.Config{Follow: true, ReOpen: true, Poll: true, CompleteLines: true, Location: seekInfo})
	if err != nil {
		logger.Err(err).Msg("error tailing file")
		return err
	}
	for {
		select {
		case <-r.Context().Done():
			return nil
		case <-s.ctx.Done():
			return nil
		case line := <-tail.Lines:
			if line != nil {
				if line.Err != nil {
					logger.Err(err).Msg("error tailing file")
					return line.Err
				}
				err := r.Send(&pb.GetJobExecutionsReply{
					Result: []string{line.Text},
				})
				if err != nil {
					logger.Err(err).Msg("error sending content")
				}
			}
		}
	}
}
