package recorder

import (
	"context"

	pb "github.com/andrescosta/jobico/internal/api/types"
)

type ExecutionRecorder interface {
	OldRecords(n int) ([]string, error)
	StartTailing(context.Context) (Tailer, error)
	AddExecution(ex *pb.JobExecution) error
	Close() error
}
type Tailer interface {
	Lines() <-chan string
	Stop() error
}
