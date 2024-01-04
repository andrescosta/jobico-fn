package recorder

import (
	"github.com/andrescosta/goico/pkg/ioutil"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogRecorder struct {
	logger zerolog.Logger
}

func New(fullpath string) (*LogRecorder, error) {
	writer := &lumberjack.Logger{
		Filename:   fullpath,
		MaxBackups: 1,
		MaxSize:    1,
		MaxAge:     1,
	}
	logger := zerolog.New(writer).With().Timestamp().Logger()
	logger.Level(zerolog.InfoLevel)
	// we create the file if not exists because tail has issues when the file is not present
	if err := ioutil.CreateEmptyIfNotExists(fullpath); err != nil {
		return nil, err
	}
	return &LogRecorder{
		logger: logger,
	}, nil
}

func (r *LogRecorder) AddExecution(ex *pb.JobExecution) error {
	r.logger.Info().
		Str("Type", ex.Result.TypeDesc).
		Str("Event", ex.Event).
		Str("Queue", ex.Queue).
		Uint64("Code", ex.Result.Code).
		Str("Result", ex.Result.Message).
		Send()
	return nil
}
