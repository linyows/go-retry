package main

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/mattn/go-shellwords"
)

// ExecCmd returns exit code
func ExecCmd(c []string) int {
	var cmd *exec.Cmd

	if len(c) > 1 {
		cmd = exec.Command(c[0], c[1:]...)
	} else {
		cmd = exec.Command(c[0])
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
		return ExitCodeError
	}

	return ExitCodeOK
}

// BuildShellCmd returns args as exec.Command
func BuildShellCmd(args []string) ([]string, error) {
	shell := os.Getenv("SHELL")
	cmd := append([]string{shell, "-c"}, args...)

	s := shellwords.NewParser()
	s.ParseEnv = true
	s.ParseBacktick = true

	return s.Parse(strings.Join(cmd, " "))
}
