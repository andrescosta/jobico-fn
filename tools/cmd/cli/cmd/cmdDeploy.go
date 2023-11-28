package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/andrescosta/workflew/api/pkg/remote"
	"github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/tools/internal/yaml"
)

var cmdDeploy = &command{

	usageLine: `cli deploy < deployment file >.yaml`,
	short:     "deploy the specified by the deployment file ",
	long:      `Deploy the file`,
}

func initDeploy() {
	cmdDeploy.flag = *flag.NewFlagSet("deploy", flag.ContinueOnError)
	cmdDeploy.run = runDeploy
	cmdDeploy.flag.Usage = func() {}

}

func runDeploy(ctx context.Context, cmd *command, args []string) {
	if len(args) < 1 {
		printHelp(os.Stdout, cmd)
		return
	}
	file := args[0]
	f, err := yaml.Decode(file)
	if err != nil {
		printError(os.Stderr, cmd, err)
		return
	} else {
		fmt.Printf("%#v\n", f.String())
		s, err := yaml.Encode(f)

		if err != nil {
			printError(os.Stderr, cmd, err)
			return
		}
		fmt.Println(*s)

		c := remote.NewControlClient()
		c.AddTenant(context.Background(), &types.Tenant{TenantId: f.TenantId})
		r, err := c.AddPackage(context.Background(), f)
		if err != nil {
			printError(os.Stderr, cmd, err)
			return
		}
		fmt.Println(r.String())
	}
}
