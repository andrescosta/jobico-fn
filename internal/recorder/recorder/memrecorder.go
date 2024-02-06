package recorder

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"

	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/rs/zerolog"
)

type MemRecorder struct {
	mux     sync.RWMutex
	results *bytes.Buffer
	tail    *MemTailer
	logger  *zerolog.Logger
}

type MemTailer struct {
	m *MemRecorder
	c chan string
}

func NewMemrecorder() *MemRecorder {
	results := bytes.NewBuffer(make([]byte, 0))
	logger := zerolog.New(results).With().Timestamp().Logger()
	return &MemRecorder{
		results: results,
		logger:  &logger,
		mux:     sync.RWMutex{},
	}
}

func (m *MemRecorder) OldRecords(n int) ([]string, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	l := make([]string, 0)
	for n != 0 {
		str, err := m.results.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err
			}
			break
		}
		l = append(l, str)
	}
	return l, nil
}

func (m *MemRecorder) Close() error {
	if m.tail != nil {
		close(m.tail.c)
	}
	return nil
}

func (m *MemRecorder) StartTailing(_ context.Context) (Tailer, error) {
	mt := &MemTailer{
		m: m,
		c: make(chan string),
	}
	m.tail = mt
	return mt, nil
}

func (t *MemTailer) Lines() <-chan string {
	return t.c
}

func (t *MemTailer) Stop() error {
	close(t.c)
	return nil
}

func (m *MemRecorder) AddExecution(ex *pb.JobExecution) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.logger.Info().
		Str("Type", ex.Result.TypeDesc).
		Str("Event", ex.Event).
		Str("Queue", ex.Queue).
		Uint64("Code", ex.Result.Code).
		Str("Result", ex.Result.Message).
		Send()
	return nil
}
