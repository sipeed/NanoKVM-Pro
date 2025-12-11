package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ShellResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// RunShell executes shell commands and returns the result
func RunShell(commands ...string) (*ShellResult, error) {
	if len(commands) == 0 {
		return &ShellResult{
			Stderr:   "no command provided",
			ExitCode: -1,
		}, fmt.Errorf("no command provided")
	}

	fullCmd := strings.Join(commands, " && ")
	cmd := exec.Command("/bin/sh", "-c", fullCmd)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	cmd.Dir = "/"
	cmd.Env = os.Environ()

	err := cmd.Run()

	result := &ShellResult{
		Stdout:   strings.TrimSpace(stdoutBuf.String()),
		Stderr:   strings.TrimSpace(stderrBuf.String()),
		ExitCode: 0,
	}

	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			result.ExitCode = exitError.ExitCode()
		}
		return result, err
	}

	return result, nil
}
