// Package config provides configuration and path generation for shotty.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config holds runtime configuration for shotty.
type Config struct {
	Hostname      string
	ScreenshotDir string
	RecordingDir  string
	StateFile     string
}

// New creates a Config from environment defaults.
func New() *Config {
	hostname, _ := os.Hostname()
	home := os.Getenv("HOME")
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}

	return &Config{
		Hostname:      hostname,
		ScreenshotDir: filepath.Join(home, "desktop", "pictures", "screenshots"),
		RecordingDir:  filepath.Join(home, "desktop", "videos", "recordings"),
		StateFile:     filepath.Join(runtimeDir, "shotty.json"),
	}
}

// GenerateScreenshotPath returns a timestamped PNG path under the screenshot directory.
func (c *Config) GenerateScreenshotPath() string {
	dir := filepath.Join(c.ScreenshotDir, c.Hostname)
	_ = os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, time.Now().Format("2006-01-02-150405")+".png")
}

// GenerateRecordingPath returns a timestamped MP4 path under the recording directory.
func (c *Config) GenerateRecordingPath() string {
	dir := filepath.Join(c.RecordingDir, c.Hostname)
	_ = os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, time.Now().Format("2006-01-02-150405")+".mp4")
}
