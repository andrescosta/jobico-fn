package recorder

import (
	"context"
	"errors"
	"io"

	"github.com/andrescosta/goico/pkg/ioutil"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/nxadm/tail"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type FileLogRecorder struct {
	path   string
	logger zerolog.Logger
	writer *lumberjack.Logger
}

func NewFileLogRecorder(path string) (ExecutionRecorder, error) {
	writer := &lumberjack.Logger{
		Filename:   path,
		MaxBackups: 1,
		MaxSize:    1,
		MaxAge:     1,
	}
	logger := zerolog.New(writer).With().Timestamp().Logger()
	logger.Level(zerolog.InfoLevel)
	// we create the file if not exists because tail has issues when the file is not present
	if err := ioutil.Touch(path); err != nil {
		return nil, err
	}
	return &FileLogRecorder{
		logger: logger,
		path:   path,
		writer: writer,
	}, nil
}

func (l *FileLogRecorder) Close() error {
	return l.writer.Close()
}

func (l *FileLogRecorder) OldRecords(n int) ([]string, error) {
	lines, err := ioutil.LastLines(l.path, n, true, true)
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func (l *FileLogRecorder) StartTailing(ctx context.Context) (Tailer, error) {
	seekInfo := &tail.SeekInfo{
		Offset: 0,
		Whence: io.SeekEnd,
	}
	tail, err := tail.TailFile(l.path, tail.Config{Follow: true, ReOpen: true, Poll: true, CompleteLines: true, Location: seekInfo})
	if err != nil {
		return nil, err
	}
	lt := &logTail{
		tail:    tail,
		linesCh: make(chan string),
	}
	go func() {
		err := lt.startTailing(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			zerolog.Ctx(ctx).Err(err).Msg("getting file results")
		}
	}()
	return lt, nil
}

type logTail struct {
	tail    *tail.Tail
	linesCh chan string
}

func (l *logTail) Lines() <-chan string {
	return l.linesCh
}

func (l *logTail) startTailing(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-l.tail.Lines:
			if !ok {
				return errors.New("closed")
			}
			l.linesCh <- line.Text
		}
	}
}

func (l *logTail) Stop() error {
	close(l.linesCh)
	err := l.tail.Stop()
	l.tail.Cleanup()
	return err
}

func (l *FileLogRecorder) AddExecution(ex *pb.JobExecution) error {
	l.logger.Info().
		Str("Type", ex.Result.TypeDesc).
		Str("Event", ex.Event).
		Str("Queue", ex.Queue).
		Uint64("Code", ex.Result.Code).
		Str("Result", ex.Result.Message).
		Send()
	return nil
}
