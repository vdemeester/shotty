package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

var version = "dev"

func main() {
	delayFlag := &cli.IntFlag{
		Name:    "delay",
		Aliases: []string{"w"},
		Usage:   "Delay capture/recording in seconds",
		Value:   0,
	}

	cmd := &cli.Command{
		Name:    "shotty",
		Usage:   "Screenshot and recording tool for Wayland",
		Version: version,
		Commands: []*cli.Command{
			// Screenshots
			{Name: "select-clipboard", Usage: "Capture selection to clipboard", Flags: []cli.Flag{delayFlag}, Action: notImplemented},
			{Name: "select-file", Usage: "Capture selection to file", Flags: []cli.Flag{delayFlag}, Action: notImplemented},
			{Name: "select-edit", Usage: "Capture selection and edit with satty", Flags: []cli.Flag{delayFlag}, Action: notImplemented},
			{Name: "window-clipboard", Usage: "Capture focused window to clipboard", Flags: []cli.Flag{delayFlag}, Action: notImplemented},
			{Name: "window-file", Usage: "Capture focused window to file", Flags: []cli.Flag{delayFlag}, Action: notImplemented},
			{Name: "screen-clipboard", Usage: "Capture focused screen to clipboard", Flags: []cli.Flag{delayFlag}, Action: notImplemented},
			{Name: "screen-file", Usage: "Capture focused screen to file", Flags: []cli.Flag{delayFlag}, Action: notImplemented},
			// Recording
			{Name: "record-select", Usage: "Record video of selection", Flags: []cli.Flag{delayFlag}, Action: notImplemented},
			{Name: "record-screen", Usage: "Record video of screen", Flags: []cli.Flag{delayFlag}, Action: notImplemented},
			{Name: "record-stop", Usage: "Stop recording and convert to mp4", Action: notImplemented},
			{Name: "record-pause", Usage: "Pause/resume recording", Action: notImplemented},
			{Name: "record-toggle", Usage: "Toggle recording (start select if idle, stop if recording)", Flags: []cli.Flag{delayFlag}, Action: notImplemented},
			// Waybar
			{Name: "waybar-status", Usage: "Output waybar status JSON", Flags: []cli.Flag{
				&cli.BoolFlag{Name: "follow", Usage: "Continuously output on state change"},
			}, Action: notImplemented},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func notImplemented(_ context.Context, _ *cli.Command) error {
	return fmt.Errorf("not implemented")
}
