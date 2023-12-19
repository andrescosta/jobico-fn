package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
)

var cmdHelp = &command{
	name:      "help",
	usageLine: "cli help < command >",
	short:     "display help for the provided command",
	long:      `Help displays usage information`,
}

func initHelp() {
	cmdHelp.flag = *flag.NewFlagSet("help", flag.ContinueOnError)
	cmdHelp.run = runHelp
	cmdHelp.flag.Usage = func() {}
}
func runHelp(_ context.Context, _ *command, args []string) {
	if len(args) < 1 || args[0] == "help" {
		printUsage(os.Stdout, cliCommand)
	}
	for _, c := range cliCommand.commands {
		if c.Name() == args[0] {
			printHelp(os.Stdout, c)
			return
		}
	}
	fmt.Printf(`unknow "%s" help topic. run "cli help"`, args[0])
}
