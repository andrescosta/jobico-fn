package cli

import (
	"context"
	"os"

	"github.com/andrescosta/goico/pkg/service"
)

func newCli() *command {
	cliCommand := &command{
		name:      "cli",
		usageLine: "cli",
		long:      "Cli is the command line admin tool.",
	}
	cliCommand.commands = []*command{
		newHelp(cliCommand),
		newUpload(),
		newDeploy(),
		newRollback(),
		newRecorder(),
		newShow(),
		newEnv(),
	}
	cliCommand.run = runCli
	return cliCommand
}

func runCli(ctx context.Context, cliCommand *command, dialer service.GrpcDialer, _ []string) {
	if len(os.Args) < 2 {
		printUsage(os.Stdout, cliCommand)
		return
	}
	cmdFound := false
	for _, command := range cliCommand.commands {
		if command.Name() == os.Args[1] {
			cmdFound = true
			if err := command.flag.Parse(os.Args[2:]); err != nil {
				printHelp(os.Stdout, command)
				return
			}
			command.run(ctx, command, dialer, command.flag.Args())
		}
	}
	if !cmdFound {
		printUsage(os.Stdout, cliCommand)
	}
}
