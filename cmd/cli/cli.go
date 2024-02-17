package main

import (
	"log"
	"os"

	"github.com/andrescosta/goico/pkg/context"
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

	ctx, cancel := context.ForEndSignals()
	defer cancel()
	cli.RunCli(ctx, os.Args)
}
