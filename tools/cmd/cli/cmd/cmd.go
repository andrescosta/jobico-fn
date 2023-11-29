package cmd

import (
	"context"
	"flag"
	"strings"
)

type command struct {
	run       func(ctx context.Context, cmd *command, args []string)
	usageLine string
	short     string
	long      string
	flag      flag.FlagSet
	commands  []*command
}

func (c *command) LongName() string {
	name := c.usageLine
	if i := strings.Index(name, " ["); i >= 0 {
		name = name[:i]
	} else {
		if i := strings.Index(name, " <"); i >= 0 {
			name = name[:i]
		}
		if i := strings.Index(name, " <"); i >= 0 {
			name = name[:i]
		}
	}
	if name == "cli" {
		return ""
	}
	return strings.TrimPrefix(name, "cli ")
}

func (c *command) Name() string {
	name := c.LongName()
	if i := strings.LastIndex(name, " "); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// Attributes exported here for to satisfy the template engine
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
