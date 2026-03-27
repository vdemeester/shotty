package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vdemeester/shotty/internal/delay"
	"github.com/vdemeester/shotty/internal/ext"
)

// SelectClipboard captures a selected region to clipboard.
func (a *App) SelectClipboard(ctx context.Context, delaySec int) error {
	delay.Countdown(a.Config.StateFile, delaySec)

	geom, err := a.Tools.Slurp(ctx)
	if err != nil {
		return fmt.Errorf("selection cancelled: %w", err)
	}

	data, err := a.Tools.GrimRegion(ctx, geom, "")
	if err != nil {
		return fmt.Errorf("capture failed: %w", err)
	}

	if err := a.Tools.WlCopy(ctx, data, "image/png"); err != nil {
		return fmt.Errorf("clipboard copy failed: %w", err)
	}

	a.Tools.NotifySimple(ctx, "Screenshot copied to clipboard", 3000)
	return nil
}

// SelectFile captures a selected region to file.
func (a *App) SelectFile(ctx context.Context, delaySec int) error {
	delay.Countdown(a.Config.StateFile, delaySec)

	geom, err := a.Tools.Slurp(ctx)
	if err != nil {
		return fmt.Errorf("selection cancelled: %w", err)
	}

	path := a.Config.GenerateScreenshotPath()
	if _, err := a.Tools.GrimRegion(ctx, geom, path); err != nil {
		return fmt.Errorf("capture failed: %w", err)
	}

	a.handleScreenshotFileActions(path)
	return nil
}

// SelectEdit captures a selected region, opens it in satty for editing.
func (a *App) SelectEdit(ctx context.Context, delaySec int) error {
	delay.Countdown(a.Config.StateFile, delaySec)

	geom, err := a.Tools.Slurp(ctx)
	if err != nil {
		return fmt.Errorf("selection cancelled: %w", err)
	}

	tmpFile := fmt.Sprintf("/tmp/shotty-%d.png", os.Getpid())
	if _, err := a.Tools.GrimRegion(ctx, geom, tmpFile); err != nil {
		return fmt.Errorf("capture failed: %w", err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	outputFile := a.Config.GenerateScreenshotPath()
	return a.Tools.Satty(ctx, tmpFile, outputFile)
}

// WindowClipboard captures the focused window to clipboard.
func (a *App) WindowClipboard(ctx context.Context, delaySec int) error {
	delay.Countdown(a.Config.StateFile, delaySec)

	if err := a.Tools.NiriScreenshotWindow(ctx, ""); err != nil {
		return fmt.Errorf("window capture failed: %w", err)
	}

	a.Tools.NotifySimple(ctx, "Window screenshot copied to clipboard", 3000)
	return nil
}

// WindowFile captures the focused window to file.
func (a *App) WindowFile(ctx context.Context, delaySec int) error {
	delay.Countdown(a.Config.StateFile, delaySec)

	path := a.Config.GenerateScreenshotPath()
	if err := a.Tools.NiriScreenshotWindow(ctx, path); err != nil {
		return fmt.Errorf("window capture failed: %w", err)
	}

	a.handleScreenshotFileActions(path)
	return nil
}

// ScreenClipboard captures the focused screen to clipboard.
func (a *App) ScreenClipboard(ctx context.Context, delaySec int) error {
	delay.Countdown(a.Config.StateFile, delaySec)

	if err := a.Tools.NiriScreenshotScreen(ctx, ""); err != nil {
		return fmt.Errorf("screen capture failed: %w", err)
	}

	a.Tools.NotifySimple(ctx, "Screen screenshot copied to clipboard", 3000)
	return nil
}

// ScreenFile captures the focused screen to file.
func (a *App) ScreenFile(ctx context.Context, delaySec int) error {
	delay.Countdown(a.Config.StateFile, delaySec)

	path := a.Config.GenerateScreenshotPath()
	if err := a.Tools.NiriScreenshotScreen(ctx, path); err != nil {
		return fmt.Errorf("screen capture failed: %w", err)
	}

	a.handleScreenshotFileActions(path)
	return nil
}

// handleScreenshotFileActions shows a notification with post-capture actions.
func (a *App) handleScreenshotFileActions(path string) {
	ctx := context.Background()
	actions := []ext.Action{
		{ID: "copy", Label: "Copy image"},
		{ID: "copypath", Label: "Copy path"},
		{ID: "edit", Label: "Edit"},
	}

	action, err := a.Tools.Notify(ctx,
		fmt.Sprintf("Screenshot saved: %s", filepath.Base(path)),
		"", 30000, actions)
	if err != nil || action == "" {
		return
	}

	switch action {
	case "copy":
		data, err := os.ReadFile(path)
		if err == nil {
			_ = a.Tools.WlCopy(ctx, data, "image/png")
		}
	case "copypath":
		_ = a.Tools.WlCopyText(ctx, path)
	case "edit":
		_ = a.Tools.Satty(ctx, path, path)
	}
}
