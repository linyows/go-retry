package main

import (
	"os"
	"os/exec"
	"syscall"
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
		} else {
			return ExitCodeError
		}
	}

	return ExitCodeOK
}
