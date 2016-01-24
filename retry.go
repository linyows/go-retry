package main

import (
	"fmt"
	"strings"
	"time"
)

// Retry returns exit status
func (cli *CLI) Retry(cmd []string) int {
	if cli.ops.UseShell {
		var cmdParseErr error
		if cmd, cmdParseErr = buildShellCmd(cmd); cmdParseErr != nil {
			cli.err(fmt.Sprintf("%#v", cmdParseErr))
			return ExitCodeError
		}
	}

	strCmd := strings.Join(cmd, " ")
	cli.out("Command: %s", strCmd)

	exitStatus := command.run(cmd)
	cli.out("Exit status: %d", exitStatus)
	if exitStatus == ExitCodeOK {
		return ExitCodeOK
	}

	d, timeParseErr := time.ParseDuration(cli.ops.Interval)
	if timeParseErr != nil {
		cli.err("%#v", timeParseErr)
		return ExitCodeError
	}

	for i := 0; i < cli.ops.Count; i++ {
		cli.out("Retrying...")
		cli.out("Sleep %s", cli.ops.Interval)
		time.Sleep(d)

		cli.out("Command: %s", strCmd)

		exitStatus = command.run(cmd)
		cli.out("Exit status: %d", exitStatus)
		if exitStatus == ExitCodeOK {
			return ExitCodeOK
		}
	}

	return exitStatus
}
