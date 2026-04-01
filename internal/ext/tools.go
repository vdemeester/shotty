package ext

import (
	"context"
	"fmt"
	"strings"
	"syscall"
)

// Action represents a notification action button.
type Action struct {
	ID    string
	Label string
}

// Tools wraps external tool invocations with a Runner for testability.
type Tools struct {
	runner Runner
}

// NewTools creates a Tools instance with the given Runner.
func NewTools(r Runner) *Tools {
	return &Tools{runner: r}
}

// DefaultTools creates a Tools instance that executes real commands.
func DefaultTools() *Tools {
	return &Tools{runner: ExecRunner{}}
}

// Slurp prompts the user to select a screen region and returns the geometry string.
func (t *Tools) Slurp(ctx context.Context) (string, error) {
	out, err := t.runner.Output(ctx, "slurp")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// GrimRegion captures a region to file. If file is "", returns PNG data on stdout.
func (t *Tools) GrimRegion(ctx context.Context, geometry, file string) ([]byte, error) {
	args := []string{"-g", geometry}
	if file == "" {
		args = append(args, "-")
		return t.runner.Output(ctx, "grim", args...)
	}
	args = append(args, file)
	return nil, t.runner.Run(ctx, "grim", args...)
}

// WlCopy copies data to clipboard with the given MIME type.
func (t *Tools) WlCopy(ctx context.Context, data []byte, mimeType string) error {
	return t.runner.RunWithStdin(ctx, data, "wl-copy", "--type", mimeType)
}

// WlCopyText copies text to clipboard.
func (t *Tools) WlCopyText(ctx context.Context, text string) error {
	return t.runner.Run(ctx, "wl-copy", text)
}

// NiriScreenshotWindow captures the focused window.
// Always copies to clipboard. If path is non-empty, also saves to file.
func (t *Tools) NiriScreenshotWindow(ctx context.Context, path string) error {
	args := []string{"msg", "action", "screenshot-window"}
	if path == "" {
		args = append(args, "--write-to-disk", "false")
	} else {
		args = append(args, "--path", path)
	}
	return t.runner.Run(ctx, "niri", args...)
}

// NiriScreenshotScreen captures the focused screen.
// Always copies to clipboard. If path is non-empty, also saves to file.
func (t *Tools) NiriScreenshotScreen(ctx context.Context, path string) error {
	args := []string{"msg", "action", "screenshot-screen"}
	if path == "" {
		args = append(args, "--write-to-disk", "false")
	} else {
		args = append(args, "--path", path)
	}
	return t.runner.Run(ctx, "niri", args...)
}

// NotifySimple sends a simple notification without actions.
func (t *Tools) NotifySimple(ctx context.Context, summary string, timeout int) {
	_ = t.runner.Run(ctx, "notify-send",
		"--app-name", "shotty",
		"-t", fmt.Sprintf("%d", timeout),
		summary,
	)
}

// Notify sends a desktop notification with action buttons.
// Returns the selected action ID (empty if dismissed).
func (t *Tools) Notify(ctx context.Context, summary, body string, timeout int, actions []Action) (string, error) {
	args := []string{
		"--app-name", "shotty",
		"--category", "recording",
		"-t", fmt.Sprintf("%d", timeout),
	}
	for _, a := range actions {
		args = append(args, "--action", fmt.Sprintf("%s=%s", a.ID, a.Label))
	}
	args = append(args, summary)
	if body != "" {
		args = append(args, body)
	}

	out, err := t.runner.Output(ctx, "notify-send", args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Satty opens the screenshot editor.
func (t *Tools) Satty(ctx context.Context, inputFile, outputFile string) error {
	return t.runner.Run(ctx, "satty",
		"--filename", inputFile,
		"--output-filename", outputFile,
		"--copy-command", "wl-copy",
	)
}

// ConvertToMP4 converts an AVI file to MP4 using ffmpeg.
func (t *Tools) ConvertToMP4(ctx context.Context, input, output string) error {
	return t.runner.Run(ctx, "ffmpeg",
		"-i", input,
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",
		"-y",
		output,
	)
}

// StartWfRecorder starts wf-recorder and returns the PID.
func (t *Tools) StartWfRecorder(ctx context.Context, geometry, file string, audio bool, audioDevice string) (int, error) {
	args := []string{"-f", file}
	if geometry != "" {
		args = append(args, "-g", geometry)
	}
	if audio || audioDevice != "" {
		args = append(args, "-a")
		if audioDevice != "" {
			args = append(args, audioDevice)
		}
	}
	return t.runner.Start(ctx, "wf-recorder", args...)
}

// StopWfRecorder sends SIGINT to wf-recorder process.
func (t *Tools) StopWfRecorder(pid int) error {
	return t.runner.Signal(pid, syscall.SIGINT)
}

// PauseWfRecorder sends SIGUSR1 to toggle pause.
func (t *Tools) PauseWfRecorder(pid int) error {
	return t.runner.Signal(pid, syscall.SIGUSR1)
}
