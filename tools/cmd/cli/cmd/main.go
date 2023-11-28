package cmd

import "context"

func RunApp(ctx context.Context, args []string) {
	initCliCommand().
		run(ctx, nil, args)
}
