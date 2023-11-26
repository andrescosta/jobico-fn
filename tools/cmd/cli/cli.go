package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/andrescosta/goico/pkg/config"
	"github.com/andrescosta/workflew/api/pkg/remote"
	"github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/tools/internal/yaml"
)

func main() {
	config.LoadEnvVariables()

	deployCmd := flag.NewFlagSet("deploy", flag.ExitOnError)
	repoCmd := flag.NewFlagSet("repo", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'deploy' or 'repo' subcommands")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "deploy":
		deployCmd.Parse(os.Args[2:])
		args := deployCmd.Args()
		if len(args) < 1 {
			fmt.Println("expected yaml deploy file name")
			os.Exit(1)
		}
		file := args[0]
		f, err := yaml.Decode(file)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%#v\n", f.String())
			s, err := yaml.Encode(f)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(*s)
			}

			c := remote.NewControlClient()
			c.AddTenant(context.Background(), &types.Tenant{TenantId: f.TenantId})
			r, err := c.AddPackage(context.Background(), f)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(r.String())
			}
		}
	case "repo":
		repoCmd.Parse(os.Args[2:])
		args := repoCmd.Args()
		if len(args) < 3 {
			fmt.Println("expected yaml deploy file name")
			os.Exit(1)
		}
		tenant := args[0]
		name := args[1]
		file := args[2]
		f, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		c := remote.NewRepoClient()
		c.AddFile(context.Background(), tenant, name, f)
	}
	//yaml.Debug()
}
