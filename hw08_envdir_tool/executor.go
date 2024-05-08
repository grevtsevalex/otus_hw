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

	cmdName := cmd[0]
	cmdArgs := cmd[1:]

	command := exec.Command(cmdName, cmdArgs...)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	command.Run()

	return command.ProcessState.ExitCode()
}
