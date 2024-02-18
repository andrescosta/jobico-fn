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

func newRollback() *command {
	cmdRollback := &command{
		name:      "rollback",
		usageLine: `cli rollabck < deployment file >.yaml`,
		short:     "rollback a Job",
		long: `
	The rollback command eliminates a Job definition from the platform, halting associated queue executors.`,
	}
	cmdRollback.flag = *flag.NewFlagSet("rollback", flag.ContinueOnError)
	cmdRollback.run = runRollback
	cmdRollback.flag.Usage = func() {}
	return cmdRollback
}

func runRollback(ctx context.Context, cmd *command, d service.GrpcDialer, args []string) {
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
	f := pb.JobPackage{}
	if err = yamlutil.DecodeFile(file, &f); err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	client, err := client.NewCtl(ctx, d)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	p, err := client.Package(ctx, f.Tenant, &f.ID)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	if len(p) < 1 {
		printError(os.Stderr, cmd, fmt.Errorf("package %s does not exist", f.ID))
		return
	}
	err = client.DeletePackage(context.Background(), &f)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	fmt.Println("The package was deleted.")
}
