package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrescosta/goico/pkg/config"
	"github.com/andrescosta/jobico/tools/internal/cmd"
)

func main() {
	config.LoadEnvVariables()
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		done()
		if r := recover(); r != nil {
			fmt.Printf("%v\n", r)
		}
	}()
	cmd.RunCli(ctx, os.Args)
	done()
}
