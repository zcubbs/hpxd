package cmd

import "os/exec"

func RunCmd(cmd string, args ...string) error {
	return exec.Command(cmd, args...).Run()
}

func RunCmdOutput(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).Output()
}

func RunCmdCombinedOutput(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).CombinedOutput()
}
