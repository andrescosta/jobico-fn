package main

import (
	"github.com/andrescosta/workflew/tools/internal/tapp"
)

func main() {
	tapp, err := tapp.New()
	if err != nil {
		panic(err)
	}
	defer tapp.Dispose()
	if err := tapp.Run(); err != nil {
		panic(err)
	}
}
