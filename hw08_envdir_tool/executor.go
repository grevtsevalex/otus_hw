package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envName, e := range env {
		if e.NeedRemove {
			os.Unsetenv(envName)
			continue
		}

		os.Setenv(envName, e.Value)
	}

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	command.Run()

	return command.ProcessState.ExitCode()
}
