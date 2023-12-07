package main

import (
	"context"
	"flag"
	"log"

	"github.com/andrescosta/jobico/tools/internal/tapp"
)

func main() {
	debugFlag := flag.Bool("debug", false, "debug enabled")
	syncUpdatesFlag := flag.Bool("sync", false, "sync enabled")
	flag.Parse()
	tapp, err := tapp.New(context.Background(), *syncUpdatesFlag)
	if err != nil {
		log.Panic(err)
	}
	if *debugFlag {
		tapp.DebugOn()
	}
	defer tapp.Dispose()
	if err := tapp.Run(); err != nil {
		panic(err)
	}
}
