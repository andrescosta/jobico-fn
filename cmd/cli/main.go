package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrescosta/goico/pkg/env"
	"github.com/andrescosta/jobico/internal/cli"
)

func main() {
	loaded, _, err := env.Load("cli")
	if err != nil {
		log.Fatalf("Error initializing %v\n", err)
	}
	if !loaded {
		log.Fatal(".env files were not loaded")
	}

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	cli.RunCli(ctx, os.Args)
	defer done()
}
