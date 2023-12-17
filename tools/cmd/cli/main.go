package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/jobico/tools/internal/cmd"
)

func main() {
	if err := env.Populate(); err != nil {
		log.Fatalf("Error initializing %v\n", err)
	}

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	cmd.RunCli(ctx, os.Args)
	defer done()
}
