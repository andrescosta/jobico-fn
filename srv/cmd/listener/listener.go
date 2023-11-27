package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/andrescosta/goico/pkg/server"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/workflew/srv/internal/listener"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func main() {
	service.StartNamed("listener", serviceFunc)
}

func serviceFunc(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	srv, err := server.New(os.Getenv("listener.addr"))
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	controler, err := listener.New(ctx)
	if err != nil {
		return err
	}

	router.Mount("/events", controler.Routes(*logger))
	logger.Info().Msgf("Events listener listening at %s", srv)
	err = srv.ServeHTTP(ctx, &http.Server{Handler: router})
	logger.Info().Msg("Events listener stopped")
	return err
}
