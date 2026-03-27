package ext

import (
	"context"
	"testing"
)

// fakeRunner records commands for verification without executing anything.
type fakeRunner struct {
	calls  []fakeCall
	output []byte // returned by Output()
	err    error  // returned by Run()/Output()/Start()
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
	return 12345, f.err
}

// --- Slurp tests ---

func TestSlurp(t *testing.T) {
	r := &fakeRunner{output: []byte("100,200 300x400\n")}
	tools := NewTools(r)

	geom, err := tools.Slurp(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if geom != "100,200 300x400" {
		t.Errorf("got %q, want %q", geom, "100,200 300x400")
	}
	if len(r.calls) != 1 || r.calls[0].name != "slurp" {
		t.Errorf("expected slurp call, got %+v", r.calls)
	}
}

// --- Grim tests ---

func TestGrimRegionToStdout(t *testing.T) {
	pngData := []byte{0x89, 0x50, 0x4e, 0x47} // fake PNG header
	r := &fakeRunner{output: pngData}
	tools := NewTools(r)

	data, err := tools.GrimRegion(context.Background(), "100,200 300x400", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 4 {
		t.Errorf("expected 4 bytes, got %d", len(data))
	}
	call := r.calls[0]
	if call.name != "grim" {
		t.Errorf("expected grim, got %s", call.name)
	}
	// Should use "-" for stdout
	assertContains(t, call.args, "-")
	assertContains(t, call.args, "-g")
}

func TestGrimRegionToFile(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	_, err := tools.GrimRegion(context.Background(), "100,200 300x400", "/tmp/test.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	call := r.calls[0]
	assertContains(t, call.args, "/tmp/test.png")
	assertNotContains(t, call.args, "-")
}

// --- Clipboard tests ---

func TestWlCopy(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	data := []byte("image data")
	err := tools.WlCopy(context.Background(), data, "image/png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	call := r.calls[0]
	if call.name != "wl-copy" {
		t.Errorf("expected wl-copy, got %s", call.name)
	}
	assertContains(t, call.args, "--type")
	assertContains(t, call.args, "image/png")
	if string(call.stdin) != "image data" {
		t.Errorf("stdin: got %q, want %q", call.stdin, "image data")
	}
}

func TestWlCopyText(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	err := tools.WlCopyText(context.Background(), "/tmp/file.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	call := r.calls[0]
	if call.name != "wl-copy" {
		t.Errorf("expected wl-copy, got %s", call.name)
	}
	assertContains(t, call.args, "/tmp/file.png")
}

// --- Niri tests ---

func TestNiriScreenshotWindowClipboardOnly(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	err := tools.NiriScreenshotWindow(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	call := r.calls[0]
	if call.name != "niri" {
		t.Errorf("expected niri, got %s", call.name)
	}
	assertContains(t, call.args, "screenshot-window")
	assertContains(t, call.args, "--write-to-disk")
	assertContains(t, call.args, "false")
}

func TestNiriScreenshotWindowToFile(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	err := tools.NiriScreenshotWindow(context.Background(), "/tmp/shot.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	call := r.calls[0]
	assertContains(t, call.args, "--path")
	assertContains(t, call.args, "/tmp/shot.png")
	assertNotContains(t, call.args, "--write-to-disk")
}

func TestNiriScreenshotScreenClipboardOnly(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	err := tools.NiriScreenshotScreen(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	call := r.calls[0]
	assertContains(t, call.args, "screenshot-screen")
	assertContains(t, call.args, "--write-to-disk")
}

func TestNiriScreenshotScreenToFile(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	err := tools.NiriScreenshotScreen(context.Background(), "/tmp/screen.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	call := r.calls[0]
	assertContains(t, call.args, "--path")
	assertContains(t, call.args, "/tmp/screen.png")
}

// --- Notify tests ---

func TestNotifySimple(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	tools.NotifySimple(context.Background(), "Screenshot taken", 3000)

	call := r.calls[0]
	if call.name != "notify-send" {
		t.Errorf("expected notify-send, got %s", call.name)
	}
	assertContains(t, call.args, "--app-name")
	assertContains(t, call.args, "shotty")
	assertContains(t, call.args, "Screenshot taken")
	assertContains(t, call.args, "3000")
}

func TestNotifyWithActions(t *testing.T) {
	r := &fakeRunner{output: []byte("copy\n")}
	tools := NewTools(r)

	actions := []Action{
		{ID: "copy", Label: "Copy image"},
		{ID: "edit", Label: "Edit"},
	}
	result, err := tools.Notify(context.Background(), "Saved", "body text", 30000, actions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "copy" {
		t.Errorf("got %q, want %q", result, "copy")
	}
	call := r.calls[0]
	assertContains(t, call.args, "--action")
	assertContains(t, call.args, "Saved")
	assertContains(t, call.args, "body text")
}

// --- Satty tests ---

func TestSatty(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	err := tools.Satty(context.Background(), "/tmp/in.png", "/tmp/out.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	call := r.calls[0]
	if call.name != "satty" {
		t.Errorf("expected satty, got %s", call.name)
	}
	assertContains(t, call.args, "--filename")
	assertContains(t, call.args, "/tmp/in.png")
	assertContains(t, call.args, "--output-filename")
	assertContains(t, call.args, "/tmp/out.png")
}

// --- FFmpeg tests ---

func TestConvertToMP4(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	err := tools.ConvertToMP4(context.Background(), "/tmp/in.avi", "/tmp/out.mp4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	call := r.calls[0]
	if call.name != "ffmpeg" {
		t.Errorf("expected ffmpeg, got %s", call.name)
	}
	assertContains(t, call.args, "-i")
	assertContains(t, call.args, "/tmp/in.avi")
	assertContains(t, call.args, "/tmp/out.mp4")
}

// --- Recorder tests ---

func TestStartWfRecorderWithGeometry(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	pid, err := tools.StartWfRecorder(context.Background(), "100,200 300x400", "/tmp/rec.avi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pid != 12345 {
		t.Errorf("pid: got %d, want 12345", pid)
	}
	call := r.calls[0]
	if call.name != "wf-recorder" {
		t.Errorf("expected wf-recorder, got %s", call.name)
	}
	assertContains(t, call.args, "-g")
	assertContains(t, call.args, "100,200 300x400")
	assertContains(t, call.args, "-f")
	assertContains(t, call.args, "/tmp/rec.avi")
}

func TestStartWfRecorderFullscreen(t *testing.T) {
	r := &fakeRunner{}
	tools := NewTools(r)

	_, err := tools.StartWfRecorder(context.Background(), "", "/tmp/rec.avi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	call := r.calls[0]
	assertNotContains(t, call.args, "-g")
}

// --- helpers ---

func assertContains(t *testing.T, args []string, want string) {
	t.Helper()
	for _, a := range args {
		if a == want {
			return
		}
	}
	t.Errorf("args %v does not contain %q", args, want)
}

func assertNotContains(t *testing.T, args []string, unwanted string) {
	t.Helper()
	for _, a := range args {
		if a == unwanted {
			t.Errorf("args %v should not contain %q", args, unwanted)
			return
		}
	}
}
