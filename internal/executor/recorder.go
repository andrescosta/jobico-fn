package executor

import (
	"context"
	"os"
	"time"

	pb "github.com/andrescosta/jobico/internal/api/types"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type recorder struct {
	cli    *cli
	tenant string
	event  string
}

func (r *recorder) sendLog(ctx context.Context, _ uint32, lvl uint32, msg string) error {
	now := time.Now()
	host, err := os.Hostname()
	if err != nil {
		host = "<error>"
	}

	return r.cli.recorder.AddJobExecution(ctx, &pb.JobExecution{
		Event:  r.event,
		Tenant: r.tenant,
		Queue:  "",
		Date: &timestamppb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
		Server: host,
		Result: &pb.JobResult{
			Type:     pb.JobResult_Log,
			TypeDesc: "log",
			Code:     uint64(lvl),
			Message:  msg,
		},
	})
}

func (r *recorder) sendResult(ctx context.Context, queue string, code uint64, result string) error {
	now := time.Now()
	host, err := os.Hostname()
	if err != nil {
		host = "<error>"
	}
	ex := &pb.JobExecution{
		Event:  r.event,
		Tenant: r.tenant,
		Queue:  queue,
		Date: &timestamppb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
		Server: host,
		Result: &pb.JobResult{
			Type:     pb.JobResult_Result,
			TypeDesc: "result",
			Code:     code,
			Message:  result,
		},
	}
	return r.cli.recorder.AddJobExecution(ctx, ex)
}
