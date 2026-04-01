package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vdemeester/shotty/internal/delay"
	"github.com/vdemeester/shotty/internal/ext"
	"github.com/vdemeester/shotty/internal/state"
)

// RecordSelect prompts for region selection then starts recording.
func (a *App) RecordSelect(ctx context.Context, delaySec int, audio bool, audioDevice string) error {
	geom, err := a.Tools.Slurp(ctx)
	if err != nil {
		return fmt.Errorf("selection cancelled: %w", err)
	}

	delay.Countdown(a.Config.StateFile, delaySec)
	return a.startRecording(ctx, geom, audio, audioDevice)
}

// RecordScreen starts recording the full screen.
func (a *App) RecordScreen(ctx context.Context, delaySec int, audio bool, audioDevice string) error {
	delay.Countdown(a.Config.StateFile, delaySec)
	return a.startRecording(ctx, "", audio, audioDevice)
}

func (a *App) startRecording(ctx context.Context, geometry string, audio bool, audioDevice string) error {
	path := a.Config.GenerateRecordingPath()

	pid, err := a.Tools.StartWfRecorder(ctx, geometry, path, audio, audioDevice)
	if err != nil {
		return fmt.Errorf("failed to start recording: %w", err)
	}

	s := &state.State{
		Recording: true,
		PID:       pid,
		File:      path,
		StartedAt: time.Now(),
	}

	if err := state.Write(a.Config.StateFile, s); err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}

	a.Tools.NotifySimple(ctx, "Recording started", 2000)
	return nil
}

// RecordStop stops the current recording.
func (a *App) RecordStop(ctx context.Context) error {
	s, err := state.Read(a.Config.StateFile)
	if err != nil {
		return fmt.Errorf("failed to read state: %w", err)
	}

	if !s.Recording {
		return fmt.Errorf("no recording in progress")
	}

	if err := a.Tools.StopWfRecorder(s.PID); err != nil {
		return fmt.Errorf("failed to stop recording: %w", err)
	}

	// Wait for wf-recorder to flush and close the file
	time.Sleep(500 * time.Millisecond)

	_ = state.Clear(a.Config.StateFile)

	// Post-recording notification actions
	actions := []ext.Action{
		{ID: "copypath", Label: "Copy path"},
		{ID: "delete", Label: "Delete"},
	}

	action, err := a.Tools.Notify(ctx,
		fmt.Sprintf("Recording saved: %s", filepath.Base(s.File)),
		"", 30000, actions)
	if err != nil || action == "" {
		return nil
	}

	switch action {
	case "copypath":
		_ = a.Tools.WlCopyText(ctx, s.File)
	case "delete":
		_ = os.Remove(s.File)
	}

	return nil
}

// RecordPause toggles pause on the current recording.
func (a *App) RecordPause(ctx context.Context) error {
	s, err := state.Read(a.Config.StateFile)
	if err != nil {
		return fmt.Errorf("failed to read state: %w", err)
	}

	if !s.Recording {
		return fmt.Errorf("no recording in progress")
	}

	if err := a.Tools.PauseWfRecorder(s.PID); err != nil {
		return fmt.Errorf("failed to pause: %w", err)
	}

	s.Paused = !s.Paused
	if err := state.Write(a.Config.StateFile, s); err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}

	if s.Paused {
		a.Tools.NotifySimple(ctx, "Recording paused", 2000)
	} else {
		a.Tools.NotifySimple(ctx, "Recording resumed", 2000)
	}

	return nil
}

// RecordToggle starts a region recording if idle, stops if recording.
func (a *App) RecordToggle(ctx context.Context, delaySec int, audio bool, audioDevice string) error {
	s, err := state.Read(a.Config.StateFile)
	if err != nil {
		return fmt.Errorf("failed to read state: %w", err)
	}

	if s.Recording {
		return a.RecordStop(ctx)
	}

	return a.RecordSelect(ctx, delaySec, audio, audioDevice)
}
