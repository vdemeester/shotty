package cmd

import (
	"context"
	"testing"
	"time"

	"github.com/vdemeester/shotty/internal/ext"
	"github.com/vdemeester/shotty/internal/state"
)

func TestRecordSelect(t *testing.T) {
	r := &fakeRunner{output: []byte("100,200 300x400\n")}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.RecordSelect(context.Background(), 0, false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have called slurp then wf-recorder
	slurpCalls := r.findCalls("slurp")
	if len(slurpCalls) != 1 {
		t.Errorf("expected 1 slurp call, got %d", len(slurpCalls))
	}

	wfCalls := r.findCalls("wf-recorder")
	if len(wfCalls) != 1 {
		t.Fatalf("expected 1 wf-recorder call, got %d", len(wfCalls))
	}
	assertArgsContain(t, wfCalls[0].args, "-g")

	// State should be written
	s, err := state.Read(cfg.StateFile)
	if err != nil {
		t.Fatalf("state read failed: %v", err)
	}
	if !s.Recording {
		t.Error("expected Recording=true in state")
	}
	if s.PID != fakePID {
		t.Errorf("PID: got %d, want %d", s.PID, fakePID)
	}
}

func TestRecordScreen(t *testing.T) {
	r := &fakeRunner{}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.RecordScreen(context.Background(), 0, false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wfCalls := r.findCalls("wf-recorder")
	if len(wfCalls) != 1 {
		t.Fatalf("expected 1 wf-recorder call, got %d", len(wfCalls))
	}
	// Fullscreen: no -g flag
	for _, a := range wfCalls[0].args {
		if a == "-g" {
			t.Error("record-screen should not pass -g flag")
		}
	}

	s, _ := state.Read(cfg.StateFile)
	if !s.Recording {
		t.Error("expected Recording=true in state")
	}
}

func TestRecordStopNoRecording(t *testing.T) {
	r := &fakeRunner{}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.RecordStop(context.Background())
	if err == nil {
		t.Fatal("expected error when no recording in progress")
	}
}

func TestRecordPauseNoRecording(t *testing.T) {
	r := &fakeRunner{}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.RecordPause(context.Background())
	if err == nil {
		t.Fatal("expected error when no recording in progress")
	}
}

func TestRecordPauseTogglesState(t *testing.T) {
	r := &fakeRunner{}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	// Write a recording state with our own PID (so it's alive)
	_ = state.Write(cfg.StateFile, &state.State{
		Recording: true,
		Paused:    false,
		PID:       fakePID,
		File:      "/tmp/test.avi",
		StartedAt: time.Now(),
	})

	err := app.RecordPause(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s, _ := state.Read(cfg.StateFile)
	if !s.Paused {
		t.Error("expected Paused=true after first pause")
	}

	// Pause again to resume
	err = app.RecordPause(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s, _ = state.Read(cfg.StateFile)
	if s.Paused {
		t.Error("expected Paused=false after second pause")
	}
}

func TestRecordToggleStartsWhenIdle(t *testing.T) {
	r := &fakeRunner{output: []byte("100,200 300x400\n")}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.RecordToggle(context.Background(), 0, false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have started recording (slurp + wf-recorder)
	wfCalls := r.findCalls("wf-recorder")
	if len(wfCalls) != 1 {
		t.Errorf("expected 1 wf-recorder call, got %d", len(wfCalls))
	}
}
