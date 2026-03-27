// Package ext provides wrappers around external tools used by shotty.
//
// All tool invocations go through the Runner interface, making them
// testable without actually executing external commands.
package ext

import (
	"bytes"
	"context"
	"os/exec"
	"syscall"
)

// Runner abstracts command execution for testability.
type Runner interface {
	// Run executes a command and waits for it to finish.
	Run(ctx context.Context, name string, args ...string) error
	// Output executes a command and returns its stdout.
	Output(ctx context.Context, name string, args ...string) ([]byte, error)
	// RunWithStdin executes a command with data piped to stdin.
	RunWithStdin(ctx context.Context, stdin []byte, name string, args ...string) error
	// Start launches a command in the background and returns its PID.
	Start(ctx context.Context, name string, args ...string) (int, error)
}

// ExecRunner implements Runner using os/exec.
type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, name string, args ...string) error {
	return exec.CommandContext(ctx, name, args...).Run()
}

func (ExecRunner) Output(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).Output()
}

func (ExecRunner) RunWithStdin(ctx context.Context, stdin []byte, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = bytes.NewReader(stdin)
	return cmd.Run()
}

func (ExecRunner) Start(ctx context.Context, name string, args ...string) (int, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}
