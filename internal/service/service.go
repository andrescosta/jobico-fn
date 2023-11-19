package service

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/go-chi/httplog"
	"github.com/joho/godotenv"
)

type Service func(context.Context) error

func Start(service Service) {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger := httplog.NewLogger("listener-log", httplog.Options{
		JSON: true,
	})

	ctx = logger.WithContext(ctx)

	err := godotenv.Load()
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
