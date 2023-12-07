package recorder

import (
	"github.com/andrescosta/goico/pkg/iohelper"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Recorder struct {
	logger zerolog.Logger
}

func NewRecorder(fullpath string) (*Recorder, error) {
	writer := &lumberjack.Logger{
		Filename:   fullpath,
		MaxBackups: 1,
		MaxSize:    1,
		MaxAge:     1,
	}
	logger := zerolog.New(writer).With().Timestamp().Logger()
	logger.Level(zerolog.InfoLevel)
	// we create the file if not exists because tail has issues when the file is not present
	if err := iohelper.CreateEmptyIfNotExists(fullpath); err != nil {
		return nil, err
	}
	return &Recorder{
		logger: logger,
	}, nil
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
