package recorder

import (
	"context"
	"io"
	"strings"

	"github.com/andrescosta/goico/pkg/ioutil"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/nxadm/tail"
	"github.com/rs/zerolog"
)

type Server struct {
	pb.UnimplementedRecorderServer

	recorder *Recorder

	fullpath string
}

func NewServer(fullpath string) (*Server, error) {
	r, err := NewRecorder(fullpath)

	if err != nil {
		return nil, err
	}

	return &Server{

		recorder: r,

		fullpath: fullpath,
	}, nil
}

func (s *Server) AddJobExecution(_ context.Context, r *pb.AddJobExecutionRequest) (*pb.AddJobExecutionReply, error) {
	if err := s.recorder.AddExecution(r.Execution); err != nil {
		return nil, err
	}

	return &pb.AddJobExecutionReply{}, nil
}

func (s *Server) GetJobExecutions(g *pb.GetJobExecutionsRequest, r pb.Recorder_GetJobExecutionsServer) error {
	logger := zerolog.Ctx(r.Context())
	seekInfo := &tail.SeekInfo{
		Offset: 0,
		Whence: io.SeekEnd,
	}

	if g.Lines != nil && *g.Lines > 0 {
		lines, err := ioutil.LastLines(s.fullpath, int(*g.Lines), true, true)
		if err != nil {
			logger.Warn().Msgf("error getting tail lines %s", err)
		} else {
			if len(lines) > 0 {
				if err := r.Send(&pb.GetJobExecutionsReply{
					Result: lines,
				}); err != nil {
					logger.Warn().Msgf("error sending tail lines %s", err)
				}
			}
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

		case line := <-tail.Lines:

			if line != nil && strings.TrimSpace(line.Text) != "" {
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
