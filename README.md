# shotty

Screenshot and screen recording tool for Wayland compositors, built for [niri](https://github.com/YaLTeR/niri).

Shotty provides a single binary with composable subcommands for capturing screenshots, recording video, and integrating with [waybar](https://github.com/Alexays/Waybar) — all driven by keybindings.

## Features

**Screenshots** — region, window, or full screen capture:
- Copy to clipboard or save to file (timestamped, organized by hostname)
- Edit with [satty](https://github.com/gabm/Satty) before saving
- Post-capture notification with actions: copy image, copy path, open editor

**Screen recording** — region or full screen video:
- Start/stop/pause/resume via separate commands (bind to keys)
- Toggle command for one-key record workflows
- Stale PID detection prevents ghost state

**Waybar integration** — live status module:
- Shows countdown timer, recording elapsed time, paused state
- Follow mode (`--follow`) for continuous polling
- CSS classes for styling: `idle`, `countdown`, `recording`, `paused`

**Delay** — optional countdown before capture/recording with waybar countdown display.

## Dependencies

| Tool | Purpose |
|------|---------|
| [grim](https://sr.ht/~emersion/grim/) | Screenshot capture |
| [slurp](https://github.com/emersion/slurp) | Region selection |
| [wl-clipboard](https://github.com/bugaevc/wl-clipboard) | Clipboard operations |
| [niri](https://github.com/YaLTeR/niri) | Window/screen capture via IPC |
| [wf-recorder](https://github.com/ammen99/wf-recorder) | Screen recording |
| [satty](https://github.com/gabm/Satty) | Screenshot annotation/editing |
| [libnotify](https://gitlab.gnome.org/GNOME/libnotify) | Desktop notifications (`notify-send`) |

## Installation

### From releases

Download the latest binary from [GitHub Releases](https://github.com/vdemeester/shotty/releases).

### From source

```sh
go install github.com/vdemeester/shotty@latest
```

Or build locally:

```sh
make build
./shotty --version
```

## Usage

```
shotty <command> [--delay N]
```

### Screenshots

```sh
shotty select-clipboard          # Select region → clipboard
shotty select-file               # Select region → file
shotty select-edit               # Select region → satty editor

shotty window-clipboard          # Focused window → clipboard
shotty window-file               # Focused window → file

shotty screen-clipboard          # Focused screen → clipboard
shotty screen-file               # Focused screen → file
```

### Recording

```sh
shotty record-select             # Select region and start recording
shotty record-screen             # Record full screen
shotty record-stop               # Stop recording, save as mp4
shotty record-pause              # Toggle pause/resume
shotty record-toggle             # Start if idle, stop if recording
```

### Waybar

```sh
shotty waybar-status             # Print status JSON once
shotty waybar-status --follow    # Continuously print on state change
```

Waybar module configuration:

```json
"custom/shotty": {
    "exec": "shotty waybar-status --follow",
    "return-type": "json",
    "format": "{}",
    "on-click": "shotty record-toggle"
}
```

### Delay

All capture and recording commands accept `--delay` / `-w` to add a countdown (in seconds) before the action starts. The countdown is visible in waybar.

```sh
shotty select-clipboard --delay 3
shotty record-screen -w 5
```

## File organization

Screenshots and recordings are saved with timestamps, grouped by hostname:

```
~/desktop/pictures/screenshots/<hostname>/2026-03-27-143022.png
~/desktop/videos/recordings/<hostname>/2026-03-27-143045.mp4
```

## License

[Apache-2.0](LICENSE)
