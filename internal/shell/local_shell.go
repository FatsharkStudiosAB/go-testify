package shell

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
)

type LocalShell struct {
	// if set to true, it will not only return the command output but also print it
	writeStdoutStderr	bool
	// if set to true, will write any output meant for stderr to stdout instead
	stderrToStdout		bool
	// if set to true, will wait for the execution to return. Otherwise runs as a background process
	waitForCompletion	bool
}

func NewLocalShell() *LocalShell {
	return &LocalShell{
		writeStdoutStderr:	true,
		stderrToStdout: 	true,
		waitForCompletion:  false,
	}
}

// Execute a shell command on the local computer.
//
// The output is returned as a byte buffer
func (shell LocalShell) Execute(ctx context.Context, cmd string, arg ...string) (bytes.Buffer, bytes.Buffer, error) {
	command := exec.CommandContext(ctx, cmd, arg...)
	return shell.ExecuteCommandAndWait(command)
}

func (shell LocalShell) ExecuteCommand(command *exec.Cmd) (bytes.Buffer, bytes.Buffer, error) {
	slog.Info(fmt.Sprintf("Executing %s", command.Path))
	stdoutBuf, stderrBuf := shell.bufferSetup(command)
	err := command.Start()
	return stdoutBuf, stderrBuf, err
}

func (shell LocalShell) ExecuteCommandAndWait(command *exec.Cmd) (bytes.Buffer, bytes.Buffer, error) {
	slog.Info(fmt.Sprintf("Executing and waiting for %s", command.Path))
	stdoutBuf, stderrBuf := shell.bufferSetup(command)
	err := command.Run()
	return stdoutBuf, stderrBuf, err
}

func (shell LocalShell) bufferSetup(command *exec.Cmd) (bytes.Buffer, bytes.Buffer) {
	var stdoutBuf, stderrBuf bytes.Buffer
	if shell.stderrToStdout {
		stderrBuf = stdoutBuf
	}
	command.Stdout = &stdoutBuf
	command.Stderr = &stderrBuf
	if shell.writeStdoutStderr {
		command.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		command.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	}
	return stdoutBuf, stderrBuf
}