package cmd

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"syscall"
	"testing"

	"github.com/vdemeester/shotty/internal/config"
	"github.com/vdemeester/shotty/internal/ext"
)

// fakePID is our own PID so state reads don't treat it as stale.
var fakePID = os.Getpid()

// fakeRunner for testing commands without executing external tools.
type fakeRunner struct {
	calls  []fakeCall
	output []byte
	err    error
}

type fakeCall struct {
	name  string
	args  []string
	stdin []byte
}

func (f *fakeRunner) Run(ctx context.Context, name string, args ...string) error {
	f.calls = append(f.calls, fakeCall{name: name, args: args})
	return f.err
}

func (f *fakeRunner) Output(ctx context.Context, name string, args ...string) ([]byte, error) {
	f.calls = append(f.calls, fakeCall{name: name, args: args})
	return f.output, f.err
}

func (f *fakeRunner) RunWithStdin(ctx context.Context, stdin []byte, name string, args ...string) error {
	f.calls = append(f.calls, fakeCall{name: name, args: args, stdin: stdin})
	return f.err
}

func (f *fakeRunner) Start(ctx context.Context, name string, args ...string) (int, error) {
	f.calls = append(f.calls, fakeCall{name: name, args: args})
	return fakePID, f.err
}

func (f *fakeRunner) Signal(_ int, _ syscall.Signal) error {
	return f.err
}

func (f *fakeRunner) findCalls(name string) []fakeCall {
	var result []fakeCall
	for _, c := range f.calls {
		if c.name == name {
			result = append(result, c)
		}
	}
	return result
}

func testConfig(t *testing.T) *config.Config {
	t.Helper()
	dir := t.TempDir()
	return &config.Config{
		Hostname:      "testhost",
		ScreenshotDir: filepath.Join(dir, "screenshots"),
		RecordingDir:  filepath.Join(dir, "recordings"),
		StateFile:     filepath.Join(dir, "shotty.json"),
	}
}

func TestSelectClipboard(t *testing.T) {
	r := &fakeRunner{output: []byte("100,200 300x400\n")}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.SelectClipboard(context.Background(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have called slurp, grim (stdout), wl-copy, notify-send
	slurpCalls := r.findCalls("slurp")
	if len(slurpCalls) != 1 {
		t.Errorf("expected 1 slurp call, got %d", len(slurpCalls))
	}

	grimCalls := r.findCalls("grim")
	if len(grimCalls) != 1 {
		t.Errorf("expected 1 grim call, got %d", len(grimCalls))
	}

	wlCopyCalls := r.findCalls("wl-copy")
	if len(wlCopyCalls) != 1 {
		t.Errorf("expected 1 wl-copy call, got %d", len(wlCopyCalls))
	}
}

func TestSelectFile(t *testing.T) {
	r := &fakeRunner{output: []byte("100,200 300x400\n")}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.SelectFile(context.Background(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	grimCalls := r.findCalls("grim")
	if len(grimCalls) != 1 {
		t.Errorf("expected 1 grim call, got %d", len(grimCalls))
	}

	// Grim should write to a file (not stdout "-")
	for _, a := range grimCalls[0].args {
		if a == "-" {
			t.Error("select-file should write to file, not stdout")
		}
	}

	// Should have screenshot dir created
	entries, _ := os.ReadDir(filepath.Join(cfg.ScreenshotDir, "testhost"))
	// Dir exists even if file doesn't (grim is faked)
	_ = entries
}

func TestWindowClipboard(t *testing.T) {
	r := &fakeRunner{}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.WindowClipboard(context.Background(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	niriCalls := r.findCalls("niri")
	if len(niriCalls) != 1 {
		t.Fatalf("expected 1 niri call, got %d", len(niriCalls))
	}
	assertArgsContain(t, niriCalls[0].args, "screenshot-window")
	assertArgsContain(t, niriCalls[0].args, "--write-to-disk")
}

func TestWindowFile(t *testing.T) {
	r := &fakeRunner{}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.WindowFile(context.Background(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	niriCalls := r.findCalls("niri")
	if len(niriCalls) != 1 {
		t.Fatalf("expected 1 niri call, got %d", len(niriCalls))
	}
	assertArgsContain(t, niriCalls[0].args, "screenshot-window")
	assertArgsContain(t, niriCalls[0].args, "--path")
}

func TestScreenClipboard(t *testing.T) {
	r := &fakeRunner{}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.ScreenClipboard(context.Background(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	niriCalls := r.findCalls("niri")
	if len(niriCalls) != 1 {
		t.Fatalf("expected 1 niri call, got %d", len(niriCalls))
	}
	assertArgsContain(t, niriCalls[0].args, "screenshot-screen")
	assertArgsContain(t, niriCalls[0].args, "--write-to-disk")
}

func TestScreenFile(t *testing.T) {
	r := &fakeRunner{}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.ScreenFile(context.Background(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	niriCalls := r.findCalls("niri")
	if len(niriCalls) != 1 {
		t.Fatalf("expected 1 niri call, got %d", len(niriCalls))
	}
	assertArgsContain(t, niriCalls[0].args, "screenshot-screen")
	assertArgsContain(t, niriCalls[0].args, "--path")
}

func TestSelectEdit(t *testing.T) {
	r := &fakeRunner{output: []byte("100,200 300x400\n")}
	tools := ext.NewTools(r)
	cfg := testConfig(t)
	app := &App{Tools: tools, Config: cfg}

	err := app.SelectEdit(context.Background(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sattyCalls := r.findCalls("satty")
	if len(sattyCalls) != 1 {
		t.Errorf("expected 1 satty call, got %d", len(sattyCalls))
	}
}

func assertArgsContain(t *testing.T, args []string, want string) {
	t.Helper()
	if !slices.Contains(args, want) {
		t.Errorf("args %v does not contain %q", args, want)
	}
}
