package cmd

import (
	"bufio"
	"fmt"
	"io"

	"github.com/andrescosta/goico/pkg/templatehelper"
)

var usageTemplate = `{{.Long | trim}}
	Usage:
	
		{{.UsageLine}} <command> [arguments]
	
	The commands are:
	{{range .Commands}}
		{{.Name | printf "%-11s"}} {{.Short}}{{end}}
	
	Use "cli help{{with .LongName}} {{.}}{{end}} <command>" for more information about a command.
	`

var helpTemplate = `usage: {{.UsageLine}}
	{{.Long | trim}}
	`

var errorTemplate = `
Error executing command: {{.Name | printf "%-11s"}}
Details:
`

func printUsage(w io.Writer, cmd *command) {
	bw := bufio.NewWriter(w)

	if err := templatehelper.Render(bw, usageTemplate, cmd); err != nil {
		fmt.Fprintf(bw, "error rendering template %v\n", err)
		return
	}
	if err := bw.Flush(); err != nil {
		fmt.Fprintf(bw, "error rendering template %v\n", err)
	}
}

func printHelp(w io.Writer, cmd *command) {
	bw := bufio.NewWriter(w)
	err := templatehelper.Render(bw, helpTemplate, cmd)
	if err != nil {
		fmt.Fprintf(bw, "error rendering template %v\n", err)
		return
	}
	if err := bw.Flush(); err != nil {
		fmt.Fprintf(bw, "error rendering template %v\n", err)
	}
}

func printError(w io.Writer, cmd *command, errCmd error) {
	bw := bufio.NewWriter(w)
	if err := templatehelper.Render(bw, errorTemplate, cmd); err != nil {
		fmt.Fprintf(bw, "error rendering template %v\n", err)
		return
	}
	fmt.Fprintln(bw, errCmd.Error())
	fmt.Fprintln(bw, "")
	if err := bw.Flush(); err != nil {
		fmt.Fprintf(bw, "error rendering template %v\n", err)
	}
}
