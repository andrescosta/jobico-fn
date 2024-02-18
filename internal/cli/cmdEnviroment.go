package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/ioutil"
	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/yamlutil"
	"github.com/andrescosta/jobico/internal/api/client"
	pb "github.com/andrescosta/jobico/internal/api/types"
)

func newEnv() *command {
	cmdEnv := &command{
		name:      "env",
		usageLine: `cli env <file>`,
		short:     "upload environment information",
		long: `
	Uploads environment information. This option is reserved for future usage.`,
	}
	cmdEnv.flag = *flag.NewFlagSet("env1", flag.ContinueOnError)
	_ = cmdEnv.flag.Bool("update", false, "override deployment")
	cmdEnv.run = runEnv
	cmdEnv.flag.Usage = func() {}
	return cmdEnv
}

func runEnv(ctx context.Context, cmd *command, d service.GrpcDialer, args []string) {
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
	client, err := client.NewCtl(ctx, d)
	if err != nil {
		printError(os.Stdout, cmd, err)
		return
	}
	var environ *pb.Environment
	environ, err = client.Environment(ctx)
	if err != nil {
		printError(os.Stdout, cmd, err)
		return
	}
	update, _ := cmd.flag.Lookup("update").Value.(flag.Getter).Get().(bool)
	if environ != nil && !update {
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
		return
	}
	fmt.Println("The environment was updated.")
}
