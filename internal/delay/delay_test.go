package delay

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/vdemeester/shotty/internal/state"
)

func TestCountdownWritesState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shotty.json")

	go Countdown(path, 2)

	// After a brief moment, state should show countdown
	time.Sleep(100 * time.Millisecond)
	s, err := state.Read(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if s.Countdown != 2 {
		t.Errorf("countdown: got %d, want 2", s.Countdown)
	}
}

func TestCountdownDecrementsEachSecond(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shotty.json")

	go Countdown(path, 2)

	// Check initial value
	time.Sleep(100 * time.Millisecond)
	s, _ := state.Read(path)
	if s.Countdown != 2 {
		t.Errorf("initial countdown: got %d, want 2", s.Countdown)
	}

	// After 1 second, should decrement
	time.Sleep(1 * time.Second)
	s, _ = state.Read(path)
	if s.Countdown != 1 {
		t.Errorf("after 1s countdown: got %d, want 1", s.Countdown)
	}
}

func TestCountdownZeroIsNoop(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shotty.json")
	Countdown(path, 0)

	// Should return immediately, no file created
	s, _ := state.Read(path)
	if s.Countdown != 0 {
		t.Errorf("expected no countdown, got %d", s.Countdown)
	}
}

func TestCountdownNegativeIsNoop(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shotty.json")
	Countdown(path, -3)

	s, _ := state.Read(path)
	if s.Countdown != 0 {
		t.Errorf("expected no countdown, got %d", s.Countdown)
	}
}
