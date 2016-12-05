package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

var outStream, errStream bytes.Buffer
var exampleCmd = "/usr/lib64/nagios/plugins/check_http -w 10 -c 15 -H localhost"

type OKCommand struct{}

func (o OKCommand) run(c []string) int {
	fmt.Fprintln(&outStream, strings.Join(c, " "))
	return 0
}

type NGCommand struct{}

func (o NGCommand) run(c []string) int {
	fmt.Fprintln(&errStream, "Error!")
	return 1
}

func TestNoArguments(t *testing.T) {
	outStream, errStream = *new(bytes.Buffer), *new(bytes.Buffer)
	command = OKCommand{}

	cli := &CLI{outStream: &outStream, errStream: &errStream}
	args := strings.Split("./retry", " ")

	if status := cli.Run(args); status != ExitCodeOK {
		t.Fatalf("expected %d, got %d.", ExitCodeOK, status)
	}

	expected := `
Usage: retry [options] command [arguments]

Options:
  -c, --count=2        retry count
  -i, --interval=3s    retry interval
  -l, --verbose        print verbose log
  -s, --shell          use shell
  -v, --version        print version information

Example:
  $ retry -i 5s -c 2 /usr/lib64/nagios/plugins/check_http -w 10 -c 15 -H localhost
`
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("expected %s, got %s.", expected, outStream.String())
	}
}

func TestVersion(t *testing.T) {
	outStream, errStream = *new(bytes.Buffer), *new(bytes.Buffer)
	command = OKCommand{}

	cli := &CLI{outStream: &outStream, errStream: &errStream}
	args := strings.Split("./retry --version", " ")

	if status := cli.Run(args); status != ExitCodeOK {
		t.Fatalf("expected %d, got %d.", ExitCodeError, status)
	}

	expected := fmt.Sprintf("retry version %s", Version)
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("expected %s, got %s.", expected, outStream.String())
	}
}

func TestInterval(t *testing.T) {
	outStream, errStream = *new(bytes.Buffer), *new(bytes.Buffer)
	command = NGCommand{}

	cli := &CLI{outStream: &outStream, errStream: &errStream}
	args := strings.Split("./retry --interval 5s "+exampleCmd, " ")

	start := time.Now()
	status := cli.Run(args)
	end := time.Now()

	if status == ExitCodeOK {
		t.Fatalf("expected not %d, got %d.", ExitCodeOK, status)
	}

	got := end.Sub(start).Seconds()
	expected := 5.
	if got < expected {
		t.Errorf("expected %f sec more, but %f sec", expected, got)
	}
}

func TestCount(t *testing.T) {
	outStream, errStream = *new(bytes.Buffer), *new(bytes.Buffer)
	command = NGCommand{}

	cli := &CLI{outStream: &outStream, errStream: &errStream}
	args := strings.Split("./retry --count 1 "+exampleCmd, " ")

	if status := cli.Run(args); status == ExitCodeOK {
		t.Fatalf("expected not %d, got %d.", ExitCodeOK, status)
	}

	expected := 2
	if count := strings.Count(errStream.String(), "Error!"); count != expected {
		t.Errorf("expected %q, got %q", expected, count)
	}
}

func TestShell(t *testing.T) {
	outStream, errStream = *new(bytes.Buffer), *new(bytes.Buffer)
	command = OKCommand{}

	cli := &CLI{outStream: &outStream, errStream: &errStream}
	os.Setenv("SHELL", "/bin/bash")
	args := strings.Split("./retry --shell "+exampleCmd, " ")

	if status := cli.Run(args); status != ExitCodeOK {
		t.Fatalf("expected %d, got %d", ExitCodeOK, status)
	}

	expected := "/bin/bash -c " + exampleCmd + "\n"
	if outStream.String() != expected {
		t.Fatalf("expected %q, got %q", expected, outStream.String())
	}
}

func TestNotShell(t *testing.T) {
	outStream, errStream = *new(bytes.Buffer), *new(bytes.Buffer)
	command = OKCommand{}

	cli := &CLI{outStream: &outStream, errStream: &errStream}
	args := strings.Split("./retry "+exampleCmd, " ")

	if status := cli.Run(args); status != ExitCodeOK {
		t.Fatalf("expected %d, got %d", ExitCodeOK, status)
	}

	expected := exampleCmd + "\n"
	if outStream.String() != expected {
		t.Fatalf("expected %q, got %q", expected, outStream.String())
	}
}
