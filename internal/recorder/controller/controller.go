package controller

import (
	"context"
	"io"
	"strings"

	"github.com/andrescosta/goico/pkg/ioutil"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/recorder/recorder"
	"github.com/nxadm/tail"
	"github.com/rs/zerolog"
)

type Recorder struct {
	recorder *recorder.LogRecorder
	fullpath string
}

func New(fullpath string) (*Recorder, error) {
	r, err := recorder.New(fullpath)
	if err != nil {
		return nil, err
	}
	return &Recorder{
		recorder: r,
		fullpath: fullpath,
	}, nil
}

func (s *Recorder) AddJobExecution(_ context.Context, r *pb.AddJobExecutionRequest) (*pb.Void, error) {
	if err := s.recorder.AddExecution(r.Execution); err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (s *Recorder) GetJobExecutions(ctx context.Context, g *pb.GetJobExecutionsRequest, r pb.Recorder_GetJobExecutionsServer) error {
	logger := zerolog.Ctx(ctx)
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
		case <-ctx.Done():
			return ctx.Err()
		case <-r.Context().Done():
			return r.Context().Err()
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
