package main

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/mattn/go-shellwords"
)

// execCmd returns exit code
func (cli *CLI) execCmd(c []string) int {
	var cmd *exec.Cmd

	if len(c) > 1 {
		cmd = exec.Command(c[0], c[1:]...)
	} else {
		cmd = exec.Command(c[0])
	}

	cmd.Stdout = cli.outStream
	cmd.Stderr = cli.errStream

	err := cmd.Run()

	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
		return ExitCodeError
	}

	return ExitCodeOK
}

// buildShellCmd returns args as exec.Command
func (cli *CLI) buildShellCmd(args []string) ([]string, error) {
	shell := os.Getenv("SHELL")
	cmd := append([]string{shell, "-c"}, args...)

	s := shellwords.NewParser()
	s.ParseEnv = true
	s.ParseBacktick = true

	return s.Parse(strings.Join(cmd, " "))
}
