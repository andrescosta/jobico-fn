package cmd

import (
	"context"
	"flag"
	"os"

	"github.com/andrescosta/jobico/api/pkg/remote"
)

var cmdRepo = &command{
	name:      "repo",
	usageLine: "cli repo <tenant> <file id> <file>.wasm|<file>.json",
	short:     "updloas the wasm or json schema file ",
	long:      `Repo updloas the wasm or json schema file`,
}

func initRepo() {
	cmdRepo.flag = *flag.NewFlagSet("repo", flag.ContinueOnError)
	cmdRepo.run = runRepo
	cmdRepo.flag.Usage = func() {}

}

func runRepo(_ context.Context, cmd *command, args []string) {
	if len(args) < 3 {
		printHelp(os.Stdout, cmd)
		return
	}
	tenant := args[0]
	name := args[1]
	file := args[2]
	f, err := os.Open(file)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	client, err := remote.NewRepoClient()
	if err != nil {
		return
	}
	if err = client.AddFile(context.Background(), tenant, name, f); err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
}
