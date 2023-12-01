package main

import (
	"github.com/andrescosta/workflew/tools/internal/tui"
)

func main() {
	cli, err := tui.NewApp()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	if err := cli.Run(); err != nil {
		panic(err)
	}
}
