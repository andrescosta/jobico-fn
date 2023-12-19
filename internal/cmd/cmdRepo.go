package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
)

var cmdRepo = &command{
	name:      "repo",
	usageLine: "cli repo <tenant> <file id> <file>.wasm <wasm>|<file>.json <json>",
	short:     "updloas the wasm or json schema file ",
	long:      `Repo updloas the wasm or json schema file`,
}

func initRepo() {
	cmdRepo.flag = *flag.NewFlagSet("repo", flag.ContinueOnError)
	cmdRepo.run = runRepo
	cmdRepo.flag.Usage = func() {}
}
func runRepo(ctx context.Context, cmd *command, args []string) {
	if len(args) < 4 {
		printHelp(os.Stdout, cmd)
		return
	}
	tenant := args[0]
	name := args[1]
	file := args[2]
	fileTypeStr := args[3]
	f, err := os.Open(file)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	client, err := remote.NewRepoClient(ctx)
	if err != nil {
		return
	}
	var fileType pb.File_FileType
	switch fileTypeStr {
	case "wasm":
		fileType = pb.File_Wasm
	case "json":
		fileType = pb.File_JsonSchema
	default:
		printHelp(os.Stdout, cmd)
		return
	}
	if err = client.AddFile(context.Background(), tenant, name, fileType, f); err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	fmt.Println("file deployed successfully")
}
