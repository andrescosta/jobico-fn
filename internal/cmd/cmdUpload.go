package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
)

var cmdUpload = &command{
	name:      "upload",
	usageLine: "cli upload <wasm|json> <tenant> <file id> <file name>",
	short:     "updload a wasm or json schema file",
	long: `
The 'upload' command enables the upload of a WebAssembly or JSON schema file to the file Repository. 
This file will be referenced by the Job definitions.`,
}

func initUpload() {
	cmdUpload.flag = *flag.NewFlagSet("upload", flag.ContinueOnError)
	cmdUpload.run = runUpload
	cmdUpload.flag.Usage = func() {}
}

func runUpload(ctx context.Context, cmd *command, args []string) {
	if len(args) < 4 {
		printHelp(os.Stdout, cmd)
		return
	}
	fileTypeStr := args[0]
	tenant := args[1]
	fileId := args[2]
	fullFileName := filepath.Clean(args[3])
	f, err := os.Open(fullFileName)
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
	if err = client.AddFile(context.Background(), tenant, fileId, fileType, f); err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	fmt.Println("file uploaded successfully")
}
