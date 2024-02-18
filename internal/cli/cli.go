package cli

import (
	"context"

	"github.com/andrescosta/goico/pkg/service"
)

func RunCli(ctx context.Context, args []string) {
	cli := newCli()
	cli.run(ctx, cli, service.DefaultGrpcDialer, args)
}
