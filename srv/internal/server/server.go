package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Server struct {
	ip       string
	port     string
	listener net.Listener
}

func New(port string) (*Server, error) {
	addr := fmt.Sprintf(":" + port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}

	return &Server{
		ip:       listener.Addr().(*net.TCPAddr).IP.String(),
		port:     strconv.Itoa(listener.Addr().(*net.TCPAddr).Port),
		listener: listener,
	}, nil
}

func (s *Server) ServeHTTP(ctx context.Context, srv *http.Server) error {
	logger := zerolog.Ctx(ctx)

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()

		logger.Debug().Msg("HTTP server: context closed")
		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		logger.Debug().Msg("HTTP server: shutting down")
		errCh <- srv.Shutdown(shutdownCtx)
	}()

	logger.Debug().Msgf("HTTP server: started on %s", s.Addr())
	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	logger.Debug().Msg("HTTP server: serving stopped")

	var merr *multierror.Error

	if err := <-errCh; err != nil {
		merr = multierror.Append(merr, fmt.Errorf("failed to shutdown server: %w", err))
	}
	return merr.ErrorOrNil()
}

func (s *Server) ServeGRPC(ctx context.Context, srv *grpc.Server) error {
	logger := zerolog.Ctx(ctx)

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()

		logger.Debug().Msg("GRPC Server: shutting down")
		srv.GracefulStop()
	}()

	logger.Debug().Msgf("GRPC Server: started on %s", s.Addr())
	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	logger.Debug().Msg("GRPC Server: stopped")

	select {
	case err := <-errCh:
		return fmt.Errorf("failed to shutdown: %w", err)
	default:
		return nil
	}

}
func (s *Server) Addr() string {
	return net.JoinHostPort(s.ip, s.port)
}

func (s *Server) IP() string {
	return s.ip
}

func (s *Server) Port() string {
	return s.port
}

func (s *Server) String() string {
	return s.Addr()
}
