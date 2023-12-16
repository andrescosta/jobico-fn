package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/ioutil"
	"github.com/andrescosta/goico/pkg/yamlutil"
	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
)

var cmdEnv = &command{

	name: "env",

	usageLine: `cli env <file>`,

	short: "uploads environment information",

	long: `Uploads env information`,
}

var cmdEnvflagUpdate *bool

func initEnv() {
	cmdEnv.flag = *flag.NewFlagSet("env1", flag.ContinueOnError)

	cmdEnvflagUpdate = cmdEnv.flag.Bool("update", false, "override deployment")

	cmdEnv.run = runEnv

	cmdEnv.flag.Usage = func() {}
}

func runEnv(ctx context.Context, cmd *command, args []string) {
	if len(args) < 1 {
		printHelp(os.Stdout, cmd)

		return
	}

	file := args[0]

	e, err := ioutil.FileExists(file)

	if err != nil {
		printError(os.Stdout, cmd, err)

		return
	}

	if !e {
		fmt.Printf("file %s does not exist.", file)

		return
	}

	client, err := remote.NewControlClient(ctx)

	if err != nil {
		printError(os.Stdout, cmd, err)
	}

	var environ *pb.Environment

	environ, err = client.GetEnviroment(ctx)

	if err != nil {
		printError(os.Stdout, cmd, err)
	}

	if environ != nil && !*cmdEnvflagUpdate {
		fmt.Println("environment exists. use -update command to override.")

		return
	}

	environ = &pb.Environment{}

	if err = yamlutil.DecodeFile(file, environ); err != nil {
		printError(os.Stderr, cmd, err)

		return
	}

	_, err = client.AddEnvironment(ctx, environ)

	if err != nil {
		printError(os.Stdout, cmd, err)
	}

	fmt.Println("The environment was updated.")
}
