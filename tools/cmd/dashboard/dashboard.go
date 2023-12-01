package main

import (
	"github.com/andrescosta/workflew/tools/internal/tapp"
)

func main() {
	app, err := tapp.New()
	if err != nil {
		panic(err)
	}
	defer app.Close()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
