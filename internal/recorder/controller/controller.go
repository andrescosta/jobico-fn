package controller

import (
	"context"
	"errors"
	"strings"

	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/recorder/recorder"
	"github.com/rs/zerolog"
)

type Recorder struct {
	recorder recorder.ExecutionRecorder
}
type Option struct {
	InMemory bool
}

func New(path string, o Option) (*Recorder, error) {
	var r recorder.ExecutionRecorder
	if o.InMemory {
		r = recorder.NewMemrecorder()
	} else {
		var err error
		r, err = recorder.NewFileLogRecorder(path)
		if err != nil {
			return nil, err
		}
	}
	return &Recorder{
		recorder: r,
	}, nil
}

func (s *Recorder) Close() error {
	return s.recorder.Close()
}

func (s *Recorder) AddJobExecution(_ context.Context, r *pb.AddJobExecutionRequest) (*pb.Void, error) {
	if err := s.recorder.AddExecution(r.Execution); err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}

func (s *Recorder) OldRecords(lines int) ([]string, error) {
	return s.recorder.OldRecords(lines)
}

func (s *Recorder) GetJobExecutionsStr(ctx context.Context, g *pb.JobExecutionsRequest, r pb.Recorder_GetJobExecutionsStrServer) error {
	logger := zerolog.Ctx(ctx)
	if g.Lines != nil && *g.Lines > 0 {
		lines, err := s.recorder.OldRecords(int(*g.Lines))
		if err != nil {
			logger.Warn().Msgf("error getting old records %s", err)
		} else {
			if len(lines) > 0 {
				if err := r.Send(&pb.JobExecutionsReply{
					Result: lines,
				}); err != nil {
					logger.Warn().Msgf("error sending tail lines %s", err)
				}
			}
		}
	}
	tail, err := s.recorder.StartTailing(r.Context())
	if err != nil {
		logger.Err(err).Msg("error tailing file")
		return err
	}
	var errLoop error
loop:
	for {
		select {
		case <-ctx.Done():
			errLoop = ctx.Err()
			break loop
		case <-r.Context().Done():
			errLoop = r.Context().Err()
			break loop
		case line, ok := <-tail.Lines():
			if !ok {
				break loop
			}
			if strings.TrimSpace(line) != "" {
				err := r.Send(&pb.JobExecutionsReply{
					Result: []string{line},
				})
				if err != nil {
					logger.Err(err).Msg("error sending content")
				}
			}
		}
	}
	errs := make([]error, 0)
	if errLoop != nil {
		errs = append(errs, errLoop)
	}

	if err := tail.Stop(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
