package cli

import (
	"context"

	"github.com/andrescosta/goico/pkg/service"
)

func RunCli(ctx context.Context, args []string) {
	initCliCommand().
		run(ctx, nil, service.DefaultGrpcDialer, args)
}
