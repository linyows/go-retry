package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

var outStream, errStream bytes.Buffer

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

func TestVersion(t *testing.T) {
	outStream, errStream = *new(bytes.Buffer), *new(bytes.Buffer)
	command = OKCommand{}

	cli := &CLI{outStream: &outStream, errStream: &errStream}
	args := strings.Split("./retry -version", " ")

	if status := cli.Run(args); status != ExitCodeOK {
		t.Fatalf("expected %d, got %d.", ExitCodeOK, status)
	}

	expected := fmt.Sprintf("retry version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Fatalf("expected %s, got %s.", expected, errStream.String())
	}
}

func TestInterval(t *testing.T) {
	outStream, errStream = *new(bytes.Buffer), *new(bytes.Buffer)
	command = NGCommand{}

	cli := &CLI{outStream: &outStream, errStream: &errStream}
	args := strings.Split("./retry -interval 5s ls -las /tmp", " ")

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
	args := strings.Split("./retry -count 1 ls -las /tmp", " ")

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
	args := strings.Split("./retry -verbose -shell=true echo $HOME", " ")

	if status := cli.Run(args); status != ExitCodeOK {
		t.Fatalf("expected %d, got %d", ExitCodeOK, status)
	}

	expected := "sh -c echo"
	if !strings.Contains(outStream.String(), expected) {
		t.Fatalf("expected %q, got %q", expected, outStream.String())
	}
}

func TestNotShell(t *testing.T) {
	outStream, errStream = *new(bytes.Buffer), *new(bytes.Buffer)
	command = OKCommand{}

	cli := &CLI{outStream: &outStream, errStream: &errStream}
	args := strings.Split("./retry -verbose -shell=false echo $HOME", " ")

	if status := cli.Run(args); status != ExitCodeOK {
		t.Fatalf("expected %d, got %d", ExitCodeOK, status)
	}

	expected := "sh -c echo"
	if strings.Contains(outStream.String(), expected) {
		t.Fatalf("expected %q, got %q", expected, outStream.String())
	}
}
