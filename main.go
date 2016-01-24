package main

import "os"

func main() {
	cli := &CLI{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}

	command = RealCommand{
		outStream: cli.outStream,
		errStream: cli.errStream,
	}

	os.Exit(cli.Run(os.Args))
}
