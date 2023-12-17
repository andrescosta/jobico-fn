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

	templatehelper.Render(bw, usageTemplate, cmd)

	bw.Flush()
}

func printHelp(w io.Writer, cmd *command) {
	bw := bufio.NewWriter(w)

	templatehelper.Render(bw, helpTemplate, cmd)

	bw.Flush()
}

func printError(w io.Writer, cmd *command, err error) {
	bw := bufio.NewWriter(w)

	templatehelper.Render(bw, errorTemplate, cmd)

	bw.Flush()

	fmt.Println(err.Error())

	fmt.Println("")
}
