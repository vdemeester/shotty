// Package state manages the recording state file for cross-invocation coordination.
package state

import (
	"encoding/json"
	"os"
	"syscall"
	"time"
)

// State represents the current recording state persisted to disk.
type State struct {
	Recording bool      `json:"recording"`
	Paused    bool      `json:"paused"`
	PID       int       `json:"pid"`
	File      string    `json:"file"`
	StartedAt time.Time `json:"started_at"`
	Countdown int       `json:"countdown,omitempty"`
}

// Read loads state from the given path. Returns an idle state if the file
// doesn't exist or contains invalid JSON. Cleans up stale PIDs automatically.
func Read(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{}, nil
		}
		return nil, err
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return &State{}, nil
	}

	// Check for stale PID
	if s.Recording && s.PID > 0 {
		if !isProcessAlive(s.PID) {
			_ = os.Remove(path)
			return &State{}, nil
		}
	}

	return &s, nil
}

// Write persists state to the given path.
func Write(path string, s *State) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Clear removes the state file.
func Clear(path string) error {
	return os.Remove(path)
}

func isProcessAlive(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}
