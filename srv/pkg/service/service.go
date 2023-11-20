package service

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/andrescosta/workflew/srv/pkg/config"
	"github.com/andrescosta/workflew/srv/pkg/log"
)

type Service func(context.Context) error

func Start(service Service) {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	err := config.LoadEnvVariables()

	logger := log.NewUsingEnv()

	ctx = logger.WithContext(ctx)

	if err != nil {
		logger.Fatal().Msg("Error loading .env file")
	}

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Fatal()
		}
	}()

	err = service(ctx)
	done()

	if err != nil {
		logger.Fatal().Err(err)
	}
}
