package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/andrescosta/goico/pkg/config"
	"github.com/andrescosta/workflew/api/pkg/remote"
	"github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/tools/internal/yaml"
)

func main() {
	config.LoadEnvVariables()

	deployCmd := flag.NewFlagSet("deploy", flag.ExitOnError)
	repoCmd := flag.NewFlagSet("repo", flag.ExitOnError)
	recorderCmd := flag.NewFlagSet("recorder", flag.ExitOnError)

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
	case "recorder":
		lines := recorderCmd.Int("lines", 0, "number of lines")
		recorderCmd.Parse(os.Args[2:])
		ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer func() {
			done()
			if r := recover(); r != nil {
				println("done error")
			}
		}()
		ch := make(chan string)
		go func(mc <-chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case l := <-mc:
					fmt.Println(l)
				}
			}

		}(ch)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := remote.NewRecorderClient().GetJobExecutions(ctx, "", int32(*lines), ch)
			if err != nil {
				fmt.Printf("error getting data %s \n", err)
			}
		}()
		fmt.Printf("getting results at proc: %d \n", os.Getpid())
		wg.Wait()
		done()
		fmt.Println("command stoped.")
	default:
		fmt.Printf("illegal option %s \n", os.Args[1])
	}

	//yaml.Debug()
}
