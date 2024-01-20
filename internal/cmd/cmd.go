package cmd

import (
	"context"
	"flag"
	"strings"

	"github.com/andrescosta/goico/pkg/service"
)

type runProc func(context.Context, *command, service.GrpcDialer, []string)

type command struct {
	run       runProc
	usageLine string
	short     string
	name      string
	long      string
	flag      flag.FlagSet
	commands  []*command
}

func (c *command) LongName() string {
	return strings.TrimPrefix(c.name, "cli ")
}

// Attributes exported here for to satisfy the template engine
func (c *command) Name() string {
	return c.name
}

func (c *command) Long() string {
	return c.long
}

func (c *command) UsageLine() string {
	return c.usageLine
}

func (c *command) Commands() []*command {
	return c.commands
}

func (c *command) Short() string {
	return c.short
}
