package config

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateScreenshotPath(t *testing.T) {
	cfg := &Config{
		Hostname:      "testhost",
		ScreenshotDir: t.TempDir(),
	}

	path := cfg.GenerateScreenshotPath()

	if !strings.HasPrefix(path, cfg.ScreenshotDir+"/testhost/") {
		t.Errorf("unexpected path prefix: %s", path)
	}
	if !strings.HasSuffix(path, ".png") {
		t.Errorf("expected .png suffix: %s", path)
	}
}

func TestGenerateRecordingPath(t *testing.T) {
	cfg := &Config{
		Hostname:     "testhost",
		RecordingDir: t.TempDir(),
	}

	path := cfg.GenerateRecordingPath()

	if !strings.HasPrefix(path, cfg.RecordingDir+"/testhost/") {
		t.Errorf("unexpected path prefix: %s", path)
	}
	if !strings.HasSuffix(path, ".mp4") {
		t.Errorf("expected .mp4 suffix: %s", path)
	}
}

func TestStateFilePath(t *testing.T) {
	t.Setenv("XDG_RUNTIME_DIR", "/tmp/shotty-test-runtime")

	cfg := New()
	expected := "/tmp/shotty-test-runtime/shotty.json"
	if cfg.StateFile != expected {
		t.Errorf("got %s, want %s", cfg.StateFile, expected)
	}
}

func TestGenerateScreenshotPathCreatesDir(t *testing.T) {
	dir := t.TempDir()
	cfg := &Config{
		Hostname:      "testhost",
		ScreenshotDir: dir,
	}

	path := cfg.GenerateScreenshotPath()

	// Verify the hostname directory was created
	info, err := os.Stat(dir + "/testhost")
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory, got file")
	}
	_ = path
}
