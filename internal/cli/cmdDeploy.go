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

func newDeploy() *command {
	cmdDeploy := &command{
		name:      "deploy",
		usageLine: `cli deploy [-update] < deployment file >.yaml`,
		short:     "deploy a Job",
		long: `
	The 'deploy' command is employed to add a job definition to the Jobico platform.
	If the '-update' flag is provided and the job has already been deployed, the command will redeploy it.`,
	}
	cmdDeploy.flag = *flag.NewFlagSet("deploy", flag.ContinueOnError)
	_ = cmdDeploy.flag.Bool("update", false, "override a deployment")
	cmdDeploy.run = runDeploy
	cmdDeploy.flag.Usage = func() {}
	return cmdDeploy
}

func runDeploy(ctx context.Context, cmd *command, d service.GrpcDialer, args []string) {
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
	isUpdate := len(p) >= 1
	update, _ := cmd.flag.Lookup("update").Value.(flag.Getter).Get().(bool)
	if isUpdate && !update {
		fmt.Printf("package %s exists. use -update command to override.\n", f.ID)
		return
	}
	t, err := client.Tenant(ctx, &f.Tenant)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	if len(t) == 0 {
		_, err = client.AddTenant(context.Background(), &pb.Tenant{ID: f.Tenant})
		if err != nil {
			printError(os.Stderr, cmd, err)
			return
		}
	}
	if !isUpdate {
		_, err = client.AddPackage(context.Background(), &f)
	} else {
		err = client.UpdatePackage(context.Background(), &f)
	}
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	if !isUpdate {
		fmt.Println("The package was deployed.")
	} else {
		fmt.Println("The package was redeployed.")
	}
}
