package main

import (
	"flag"
	"log"
	"os"

	"github.com/andrescosta/jobico/tools/internal/tapp"
)

func main() {
	debugFlag := flag.Bool("debug", false, "debug enabled")
	syncUpdatesFlag := flag.Bool("sync", false, "sync enabled")
	flag.Parse()
	tapp, err := tapp.New(*syncUpdatesFlag)
	if err != nil {
		log.Panic(err)
	}
	if len(os.Args) > 1 {
		if *debugFlag {
			tapp.CollectDebugInfo()
		}
	}
	defer tapp.Dispose()
	if err := tapp.Run(); err != nil {
		panic(err)
	}
}
