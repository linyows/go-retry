package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

type CLI struct {
	outStream, errStream io.Writer
}

func (cli *CLI) Run(args []string) int {
	var (
		interval  string
		count     int
		arguments []string
		verbose   bool
		version   bool
	)

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&interval, "interval", "3s", "Retry interval")
	flags.StringVar(&interval, "i", "3s", "Retry interval(Short)")
	flags.IntVar(&count, "count", 2, "Retry count")
	flags.IntVar(&count, "c", 2, "Retry count(Short)")
	flags.BoolVar(&verbose, "verbose", false, "Print verbose log.")
	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, strings.Fields(flags.Arg(0))...)
		flags.Parse(flags.Args()[1:])
	}

	if verbose {
		log.Printf("Command: %s", strings.Join(arguments, " "))
	}

	var cmd *exec.Cmd

	if len(arguments) > 1 {
		cmd = exec.Command(arguments[0], arguments[1:]...)
	} else {
		cmd = exec.Command(arguments[0])
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		return ExitCodeOK
	}

	sum := 0
	for i := 0; i < count; i++ {
		sum += i
		if verbose {
			log.Printf("Retrying...")
		}
		d, parseerr := time.ParseDuration(interval)
		if parseerr != nil {
			log.Fatal(parseerr)
		}

		if verbose {
			log.Printf("Sleep %s", interval)
		}
		time.Sleep(d)

		if verbose {
			log.Printf("Command: %s", strings.Join(arguments, " "))
		}

		if len(arguments) > 1 {
			cmd = exec.Command(arguments[0], arguments[1:]...)
		} else {
			cmd = exec.Command(arguments[0])
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return ExitCodeOK
		}
	}

	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			if verbose {
				log.Printf("Exit status: %d", status.ExitStatus())
			}
			return status.ExitStatus()
		} else {
			return ExitCodeError
		}
	}

	return ExitCodeOK
}
