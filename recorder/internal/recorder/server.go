package recorder

import (
	"context"
	"io"

	"github.com/nxadm/tail"
	"github.com/rs/zerolog"

	pb "github.com/andrescosta/workflew/api/types"
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
					Result: line.Text,
				})
				if err != nil {
					logger.Err(err).Msg("error sending content")
				}
			}
		}
	}
}
