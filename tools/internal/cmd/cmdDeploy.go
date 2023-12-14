package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/ioutl"
	"github.com/andrescosta/goico/pkg/yamlutl"
	"github.com/andrescosta/jobico/api/pkg/remote"
	pb "github.com/andrescosta/jobico/api/types"
)

var cmdDeploy = &command{
	name:      "deploy",
	usageLine: `cli deploy < deployment file >.yaml`,
	short:     "deploy the specified by the deployment file ",
	long:      `Deploy the file`,
}

var cmdDeployflagUpdate *bool

func initDeploy() {
	cmdDeploy.flag = *flag.NewFlagSet("deploy", flag.ContinueOnError)
	cmdDeployflagUpdate = cmdDeploy.flag.Bool("update", false, "override deployment")
	cmdDeploy.run = runDeploy
	cmdDeploy.flag.Usage = func() {}

}

func runDeploy(ctx context.Context, cmd *command, args []string) {
	if len(args) < 1 {
		printHelp(os.Stdout, cmd)
		return
	}
	file := args[0]
	e, err := ioutl.FileExists(file)
	if err != nil {
		printError(os.Stdout, cmd, err)
		return
	}
	if !e {
		fmt.Printf("file %s does not exist.", file)
		return
	}
	f := pb.JobPackage{}
	if err = yamlutl.DecodeFile(file, &f); err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	client, err := remote.NewControlClient(ctx)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}

	p, err := client.GetPackage(ctx, f.TenantId, &f.ID)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	isUpdate := len(p) >= 1
	if isUpdate && !*cmdDeployflagUpdate {
		fmt.Printf("package %s exists. use -update command to override.\n", f.ID)
		return
	}

	t, err := client.GetTenant(ctx, &f.TenantId)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	}
	if len(t) == 0 {
		_, err = client.AddTenant(context.Background(), &pb.Tenant{ID: f.TenantId})
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
