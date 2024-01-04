package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/yamlutil"
	"github.com/andrescosta/jobico/api/pkg/remote"
)

var cmdShow = &command{
	name:      "show",
	usageLine: `cli show <deploy|env> <tenant id> <deploy id>`,
	short:     "print deployments and environment information",
	long: ` 
The 'show' command prints information about Job Definitions deployed on the platform as well as 
environment informent. It offers details on the configuration, logic, and associated schema 
of a deployed job.`,
}

func initShow() {
	cmdShow.flag = *flag.NewFlagSet("show", flag.ContinueOnError)
	cmdShow.run = runShow
	cmdShow.flag.Usage = func() {}
}

func runShow(ctx context.Context, cmd *command, args []string) {
	switch args[0] {
	case "deploy":
		showDeploy(ctx, args, cmd)
	case "env":
		showEnv(ctx, args, cmd)
	default:
		printHelp(os.Stdout, cmd)
		return
	}
}

func showDeploy(ctx context.Context, args []string, cmd *command) {
	if len(args) < 3 {
		printHelp(os.Stdout, cmd)
		return
	}
	tenant := args[1]
	id := args[2]
	client, err := remote.NewControlClient(ctx)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	p, err := client.GetPackage(ctx, tenant, &id)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	if len(p) == 0 {
		fmt.Println("not found")
		return
	}
	s, err := yamlutil.Marshal(p[0])
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	fmt.Println(*s)
}

func showEnv(ctx context.Context, _ []string, cmd *command) {
	client, err := remote.NewControlClient(ctx)
	if err != nil {
		return
	}
	p, err := client.GetEnvironment(ctx)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	if p == nil {
		fmt.Println("not found")
		return
	}
	s, err := yamlutil.Marshal(p)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	fmt.Println(*s)
}
