package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/iohelper"
	"github.com/andrescosta/goico/pkg/yamlico"
	"github.com/andrescosta/workflew/api/pkg/remote"
	pb "github.com/andrescosta/workflew/api/types"
)

var cmdEnv = &command{
	name:      "env",
	usageLine: `cli env <ll>`,
	short:     "display enviroment information",
	long:      `Display env information`,
}

func initEnv() {
	cmdEnv.flag = *flag.NewFlagSet("env1", flag.ContinueOnError)
	cmdEnv.run = runEnv
	cmdEnv.flag.Usage = func() {}
}

func runEnv(ctx context.Context, cmd *command, args []string) {
	if len(args) < 1 {
		printHelp(os.Stdout, cmd)
		return
	}
	file := args[0]
	e, err := iohelper.FileExists(file)
	if err != nil {
		printError(os.Stdout, cmd, err)
		return
	}
	if !e {
		fmt.Printf("file %s does not exist.", file)
		return
	}
	environ := &pb.Environment{}
	if err = yamlico.Decode(file, environ); err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	client, err := remote.NewControlClient()
	if err != nil {
		printError(os.Stdout, cmd, err)
	}
	_, err = client.AddEnvironment(ctx, environ)
	if err != nil {
		printError(os.Stdout, cmd, err)
	}
}
