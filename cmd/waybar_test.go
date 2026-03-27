package cmd

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/vdemeester/shotty/internal/state"
)

func TestBuildWaybarStatusIdle(t *testing.T) {
	s := &state.State{}
	status := buildWaybarStatus(s)

	if status.Class != "idle" {
		t.Errorf("class: got %s, want idle", status.Class)
	}
	if status.Text != "" {
		t.Errorf("text: got %q, want empty", status.Text)
	}
}

func TestBuildWaybarStatusRecording(t *testing.T) {
	s := &state.State{
		Recording: true,
		File:      "/tmp/test.avi",
		StartedAt: time.Now().Add(-2*time.Minute - 35*time.Second),
	}
	status := buildWaybarStatus(s)

	if status.Class != "recording" {
		t.Errorf("class: got %s, want recording", status.Class)
	}
	if status.Text == "" {
		t.Error("text should not be empty during recording")
	}
}

func TestBuildWaybarStatusPaused(t *testing.T) {
	s := &state.State{Recording: true, Paused: true, StartedAt: time.Now()}
	status := buildWaybarStatus(s)

	if status.Class != "paused" {
		t.Errorf("class: got %s, want paused", status.Class)
	}
}

func TestBuildWaybarStatusCountdown(t *testing.T) {
	s := &state.State{Countdown: 3}
	status := buildWaybarStatus(s)

	if status.Class != "countdown" {
		t.Errorf("class: got %s, want countdown", status.Class)
	}
	if status.Text == "" {
		t.Error("text should not be empty during countdown")
	}
}

func TestBuildWaybarStatusJSON(t *testing.T) {
	tests := []struct {
		name  string
		state *state.State
	}{
		{"idle", &state.State{}},
		{"recording", &state.State{Recording: true, StartedAt: time.Now()}},
		{"paused", &state.State{Recording: true, Paused: true, StartedAt: time.Now()}},
		{"countdown", &state.State{Countdown: 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := buildWaybarStatus(tt.state)

			data, err := json.Marshal(status)
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}

			var parsed map[string]interface{}
			if err := json.Unmarshal(data, &parsed); err != nil {
				t.Fatalf("output is not valid JSON: %v", err)
			}

			// Verify required waybar fields exist
			for _, field := range []string{"text", "tooltip", "class", "alt"} {
				if _, ok := parsed[field]; !ok {
					t.Errorf("missing field %q in JSON output", field)
				}
			}
		})
	}
}

func TestWaybarStatusOneShot(t *testing.T) {
	cfg := testConfig(t)
	app := &App{Config: cfg}

	// No state file — should produce idle status without error
	status, err := app.BuildCurrentStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Class != "idle" {
		t.Errorf("class: got %s, want idle", status.Class)
	}
}
