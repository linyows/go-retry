package main

import (
	"fmt"
	"io"

	flag "github.com/linyows/mflag"
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

var usageText = `
Usage: retry [options] command [arguments]

Options:`

var exampleText = `
Example:
  $ retry -i 5s -c 2 /usr/lib64/nagios/plugins/check_http -w 10 -c 15 -H localhost
`

// Run for retry
func (cli *CLI) Run(args []string) int {
	f := flag.NewFlagSet(Name, flag.ContinueOnError)
	f.SetOutput(cli.outStream)

	f.Usage = func() {
		fmt.Fprintf(cli.outStream, usageText)
		f.PrintDefaults()
		fmt.Fprint(cli.outStream, exampleText)
	}

	f.StringVar(&cli.ops.Interval, []string{"i", "-interval"}, "3s", "retry interval")
	f.IntVar(&cli.ops.Count, []string{"c", "-count"}, 2, "retry count")
	f.BoolVar(&cli.ops.UseShell, []string{"s", "-shell"}, false, "use shell")
	f.BoolVar(&cli.ops.Verbose, []string{"l", "-verbose"}, false, "print verbose log")
	f.BoolVar(&cli.ops.Version, []string{"v", "-version"}, false, "print version information")

	if err := f.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if cli.ops.Version {
		fmt.Fprintf(cli.outStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	if len(f.Args()) == 0 {
		f.Usage()
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
