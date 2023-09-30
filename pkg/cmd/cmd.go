// Package cmd provides utility functions to execute system commands.
//
// Author: zakaria.elbouwab
package cmd

import "os/exec"

// RunCmd executes a system command without capturing its output.
//
// The function takes the command to be run as the first argument followed by its arguments.
// For example, to run "ls -l", call RunCmd("ls", "-l").
//
// Returns an error if the command exits with a non-zero status.
func RunCmd(cmd string, args ...string) error {
	return exec.Command(cmd, args...).Run()
}

// RunCmdOutput executes a system command and returns its standard output.
func RunCmdOutput(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).Output()
}

// RunCmdCombinedOutput executes a system command and returns its combined
// standard output and standard error.
func RunCmdCombinedOutput(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).CombinedOutput()
}
