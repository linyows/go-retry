package main

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/mattn/go-shellwords"
)

// Command interface
type Command interface {
	run([]string) int
}

// RealCommand structure
type RealCommand struct {
	outStream io.Writer
	errStream io.Writer
}

// command
var command Command

// run returns exit code
func (r RealCommand) run(c []string) int {
	var cmd *exec.Cmd

	if len(c) > 1 {
		cmd = exec.Command(c[0], c[1:]...)
	} else {
		return ExitCodeError
	}

	cmd.Stdout = r.outStream
	cmd.Stderr = r.errStream

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
func buildShellCmd(args []string) ([]string, error) {
	shell := os.Getenv("SHELL")
	cmd := append([]string{shell, "-c"}, args...)

	s := shellwords.NewParser()
	s.ParseEnv = true
	s.ParseBacktick = true

	return s.Parse(strings.Join(cmd, " "))
}
