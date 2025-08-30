package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainNoArgs(t *testing.T) {
	output, exitCode := runMainWithArgs()
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(output, "no website provided") {
		t.Errorf("expected output to contain 'no website provided', got: %s", output)
	}
}

func TestMainTooManyArgs(t *testing.T) {
	output, exitCode := runMainWithArgs("one", "two")
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(output, "too many arguments provided") {
		t.Errorf("expected output to contain 'too many arguments provided', got: %s", output)
	}
}

func TestMainValidArg(t *testing.T) {
	output, exitCode := runMainWithArgs("https://example.com")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(output, "starting crawl of: https://example.com") {
		t.Errorf("unexpected output: %s", output)
	}
}

func runMainWithArgs(args ...string) (output string, exitCode int) {
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GO_TEST_SUBPROCESS=1")

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return outBuf.String(), exitErr.ExitCode()
		}
		return outBuf.String(), -1
	}
	return outBuf.String(), 0
}

func TestMain(m *testing.M) {
	if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
		main() // Run the actual main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}
