package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteAndRead(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shotty.json")

	s := &State{
		Recording: true,
		Paused:    false,
		PID:       os.Getpid(), // use our own PID so it's alive
		File:      "/tmp/test.avi",
		StartedAt: time.Now().Truncate(time.Second),
	}

	if err := Write(path, s); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if got.PID != os.Getpid() {
		t.Errorf("PID: got %d, want %d", got.PID, os.Getpid())
	}
	if !got.Recording {
		t.Error("expected Recording=true")
	}
	if got.File != "/tmp/test.avi" {
		t.Errorf("File: got %s, want /tmp/test.avi", got.File)
	}
}

func TestReadMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")

	s, err := Read(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if s.Recording {
		t.Error("expected idle state for missing file")
	}
}

func TestClear(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shotty.json")

	_ = Write(path, &State{Recording: true, PID: 1})
	if err := Clear(path); err != nil {
		t.Fatalf("clear failed: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected state file to be removed")
	}
}

func TestReadStaleProcess(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shotty.json")

	// PID 99999999 almost certainly doesn't exist
	_ = Write(path, &State{Recording: true, PID: 99999999, File: "/tmp/test.avi"})

	s, err := Read(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if s.Recording {
		t.Error("expected stale PID to be cleaned up")
	}
}

func TestReadInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shotty.json")
	_ = os.WriteFile(path, []byte("not json"), 0o600)

	s, err := Read(path)
	if err != nil {
		t.Fatalf("expected no error for invalid JSON, got: %v", err)
	}
	if s.Recording {
		t.Error("expected idle state for invalid JSON")
	}
}

func TestWriteCountdown(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shotty.json")

	s := &State{Countdown: 5}
	if err := Write(path, s); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if got.Countdown != 5 {
		t.Errorf("Countdown: got %d, want 5", got.Countdown)
	}
}
