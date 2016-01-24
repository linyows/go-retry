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
	ops                  Ops
}

// Run for retry
func (cli *CLI) Run(args []string) int {
	f := flag.NewFlagSet(Name, flag.ContinueOnError)
	f.SetOutput(cli.errStream)

	f.StringVar(&cli.ops.Interval, "interval", "3s", "Retry interval")
	f.StringVar(&cli.ops.Interval, "i", "3s", "Retry interval(Short)")
	f.IntVar(&cli.ops.Count, "count", 2, "Retry count")
	f.IntVar(&cli.ops.Count, "c", 2, "Retry count(Short)")
	f.BoolVar(&cli.ops.UseShell, "shell", false, "Use shell")
	f.BoolVar(&cli.ops.Verbose, "verbose", false, "Print verbose log.")
	f.BoolVar(&cli.ops.Version, "version", false, "Print version information and quit.")

	if err := f.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if cli.ops.Version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	return cli.Retry(f.Args())
}

// out
func (cli *CLI) out(format string, a ...interface{}) {
	if cli.ops.Verbose {
		fmt.Fprintln(cli.outStream, fmt.Sprintf(format, a...))
	}
}

// err
func (cli *CLI) err(format string, a ...interface{}) {
	fmt.Fprintln(cli.errStream, fmt.Sprintf(format, a...))
}
