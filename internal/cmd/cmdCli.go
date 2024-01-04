package cmd

import (
	"context"
	"os"
)

func initCliCommand() *command {
	initHelp()
	initDeploy()
	initRecorder()
	initUpload()
	initShow()
	initEnv()
	initRollback()
	cliCommand.commands = []*command{
		cmdHelp,
		cmdUpload,
		cmdDeploy,
		cmdRollback,
		cmdRecorder,
		cmdShow,
		cmdEnv,
	}
	cliCommand.run = runCli
	return cliCommand
}

var cliCommand = &command{
	name:      "cli",
	usageLine: "cli",
	long:      "Cli is the command line admin tool.",
}

func runCli(ctx context.Context, _ *command, _ []string) {
	if len(os.Args) < 2 {
		printUsage(os.Stdout, cliCommand)
		return
	}
	cmdFound := false
	for _, c := range cliCommand.commands {
		if c.Name() == os.Args[1] {
			cmdFound = true
			if err := c.flag.Parse(os.Args[2:]); err != nil {
				printHelp(os.Stdout, c)
				return
			}
			c.run(ctx, c, c.flag.Args())
		}
	}
	if !cmdFound {
		printUsage(os.Stdout, cliCommand)
	}
}
