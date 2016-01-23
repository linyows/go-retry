package main

import (
	"flag"
	"fmt"
	"io"
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
	f := flag.NewFlagSet(Name, flag.ContinueOnError)
	f.SetOutput(cli.errStream)

	var ops Ops
	f.StringVar(&ops.Interval, "interval", "3s", "Retry interval")
	f.StringVar(&ops.Interval, "i", "3s", "Retry interval(Short)")
	f.IntVar(&ops.Count, "count", 2, "Retry count")
	f.IntVar(&ops.Count, "c", 2, "Retry count(Short)")
	f.BoolVar(&ops.UseShell, "shell", true, "Use shell")
	f.BoolVar(&ops.UseShell, "s", true, "Use shell(Short)")
	f.BoolVar(&ops.Verbose, "verbose", false, "Print verbose log.")
	f.BoolVar(&ops.Version, "version", false, "Print version information and quit.")

	if err := f.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if ops.Version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	return Retry(f.Args(), ops)
}
