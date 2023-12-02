package cmd

import "context"

func RunCli(ctx context.Context, args []string) {
	initCliCommand().
		run(ctx, nil, args)
}
