package main

import (
	"flag"
	"fmt"
	"io"

	"strings"
)

const (
	// ExitCodeOK for success
	ExitCodeOK int = 0
	// ExitCodeError for error
	ExitCodeError int = 1 + iota
)

// CLI is structure
type CLI struct {
	outStream, errStream io.Writer
}

// Run for retry
func (cli *CLI) Run(args []string) int {
	var arguments []string

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	var ops Ops
	flags.StringVar(&ops.Interval, "interval", "3s", "Retry interval")
	flags.StringVar(&ops.Interval, "i", "3s", "Retry interval(Short)")
	flags.IntVar(&ops.Count, "count", 2, "Retry count")
	flags.IntVar(&ops.Count, "c", 2, "Retry count(Short)")
	flags.BoolVar(&ops.Verbose, "verbose", false, "Print verbose log.")
	flags.BoolVar(&ops.Version, "version", false, "Print version information and quit.")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if ops.Version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, strings.Fields(flags.Arg(0))...)
		flags.Parse(flags.Args()[1:])
	}

	return Retry(arguments, ops)
}
