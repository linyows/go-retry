package main

import (
	"log"
	"strings"
	"time"
)

// Retry returns exit status
func Retry(cmd []string, ops Ops) int {
	if ops.UseShell {
		var err error
		if cmd, err = BuildShellCmd(cmd); err != nil {
			log.Fatal(err)
		}
	}

	strCmd := strings.Join(cmd, " ")

	if ops.Verbose {
		log.Printf("Command: %s", strCmd)
	}

	exitStatus := ExecCmd(cmd)
	if ops.Verbose {
		log.Printf("Exit status: %d", exitStatus)
	}
	if exitStatus == ExitCodeOK {
		return ExitCodeOK
	}

	d, parseerr := time.ParseDuration(ops.Interval)
	if parseerr != nil {
		log.Fatal(parseerr)
	}

	for i := 0; i < ops.Count; i++ {
		if ops.Verbose {
			log.Printf("Retrying...")
			log.Printf("Sleep %s", ops.Interval)
		}
		time.Sleep(d)

		if ops.Verbose {
			log.Printf("Command: %s", strCmd)
		}

		exitStatus = ExecCmd(cmd)
		if ops.Verbose {
			log.Printf("Exit status: %d", exitStatus)
		}
		if exitStatus == ExitCodeOK {
			return ExitCodeOK
		}
	}

	return exitStatus
}
