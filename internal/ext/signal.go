package ext

import (
	"os"
	"syscall"
)

// StopWfRecorder sends SIGINT to wf-recorder process.
func StopWfRecorder(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(syscall.SIGINT)
}

// PauseWfRecorder sends SIGUSR1 to toggle pause.
func PauseWfRecorder(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(syscall.SIGUSR1)
}
