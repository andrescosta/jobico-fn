package recorder

import (
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Recorder struct {
	logger zerolog.Logger
}

func NewRecorder(fullpath string) *Recorder {
	writer := &lumberjack.Logger{
		Filename:   fullpath,
		MaxBackups: 1,
		MaxSize:    1,
		MaxAge:     1,
	}
	logger := zerolog.New(writer).With().Timestamp().Logger()
	logger.Level(zerolog.InfoLevel)

	return &Recorder{
		logger: logger,
	}
}

func (r *Recorder) AddExecution(ex *pb.JobExecution) error {
	r.logger.Info().
		Str("Event", ex.EventId).
		Str("Queue", ex.QueueId).
		Uint64("Code", ex.Result.Code).
		Str("Result", ex.Result.Message).
		Send()
	return nil
}

func (r *Recorder) GetExecutions(tenantId string, ex *pb.JobExecution) ([]*pb.JobExecution, error) {
	return nil, nil
}
