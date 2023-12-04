package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/yamlico"
	"github.com/andrescosta/jobico/api/pkg/remote"
)

var cmdShow = &command{
	name:      "show",
	usageLine: `cli show <deploy|tenant|env> <tenant id> <deploy id>`,
	short:     "display information of the specified entity ",
	long:      `deployDisplay information`,
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
	client, err := remote.NewControlClient()
	if err != nil {
		return
	}

	p, err := client.GetPackage(ctx, tenant, &id)
	if err != nil {
		printError(os.Stderr, cmd, err)
	}
	if len(p) == 0 {
		fmt.Println("not found")
	} else {
		s, err := yamlico.Encode(p[0])
		if err != nil {
			printError(os.Stderr, cmd, err)
			return
		}
		fmt.Println(*s)
	}
}

func showEnv(ctx context.Context, args []string, cmd *command) {
	client, err := remote.NewControlClient()
	if err != nil {
		return
	}
	p, err := client.GetEnviroment(ctx)
	if err != nil {
		printError(os.Stderr, cmd, err)
	}
	if p == nil {
		fmt.Println("not found")
	} else {
		s, err := yamlico.Encode(p)
		if err != nil {
			printError(os.Stderr, cmd, err)
			return
		}
		fmt.Println(*s)
	}
}
