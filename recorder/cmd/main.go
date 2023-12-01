package main

import (
	"context"
	"log"

	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/recorder/internal/recorder"
)

func main() {
	svc, err := service.NewGrpService(context.Background(), "recorder",
		&pb.Recorder_ServiceDesc,
		func(ctx context.Context) (any, error) {
			return recorder.NewServer(".\\log.log")
		}, service.EmptyhealthCheckHandler)
	if err != nil {
		log.Panicf("error starting recorder service: %s", err)
	}
	if err = svc.Serve(); err != nil {
		log.Fatalf("error serving recorder service: %s", err)
	}

}
