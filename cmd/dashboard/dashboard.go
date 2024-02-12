package main

import (
	"context"
	"flag"
	"log"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/jobico/internal/dashboard"
)

func main() {
	debugFlag := flag.Bool("debug", false, "debug enabled")

	syncUpdatesFlag := flag.Bool("sync", false, "sync enabled")

	flag.Parse()

	d, err := dashboard.New(context.Background(), service.DefaultGrpcDialer, "dashboard", *syncUpdatesFlag)
	if err != nil {
		log.Panic(err)
	}
	defer d.Dispose()
	if *debugFlag {
		d.DebugOn()
	}
	if err := d.Run(); err != nil {
		panic(err)
	}
}
