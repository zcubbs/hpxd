package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	// Test a command that should succeed
	err := RunCmd("echo", "hello")
	if err != nil {
		t.Fatalf("Expected command to succeed, but got error: %v", err)
	}

	// Test a command that should fail
	err = RunCmd("ls", "/nonexistent")
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded")
	}
}

func TestRunCmdOutput(t *testing.T) {
	// Test a command that should succeed and return output
	output, err := RunCmdOutput("echo", "hello")
	if err != nil {
		t.Fatalf("Expected command to succeed, but got error: %v", err)
	}
	if !bytes.Equal(output, []byte("hello\n")) {
		t.Fatalf("Expected output to be 'hello\\n', but got: %s", output)
	}

	// Test a command that should fail
	_, err = RunCmdOutput("ls", "/nonexistent")
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded")
	}
}

func TestRunCmdCombinedOutput(t *testing.T) {
	// Test a command that should succeed and return output
	output, err := RunCmdCombinedOutput("echo", "hello")
	if err != nil {
		t.Fatalf("Expected command to succeed, but got error: %v", err)
	}
	if !bytes.Equal(output, []byte("hello\n")) {
		t.Fatalf("Expected output to be 'hello\\n', but got: %s", output)
	}

	// Test a command that should fail and return stderr
	output, err = RunCmdCombinedOutput("ls", "/nonexistent")
	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded")
	}
	if !strings.Contains(string(output), "No such file or directory") {
		t.Fatalf("Expected error output to contain 'No such file or directory', but got: %s", output)
	}
}
