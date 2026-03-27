package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/vdemeester/shotty/cmd"
)

var version = "dev"

func main() {
	app := cmd.NewApp()

	delayFlag := &cli.IntFlag{
		Name:    "delay",
		Aliases: []string{"w"},
		Usage:   "Delay capture/recording in seconds",
		Value:   0,
	}

	root := &cli.Command{
		Name:                   "shotty",
		Usage:                  "Screenshot and recording tool for Wayland",
		Version:                version,
		EnableShellCompletion:  true,
		Commands: []*cli.Command{
			// Screenshots
			{
				Name: "select-clipboard", Usage: "Capture selection to clipboard",
				Flags:  []cli.Flag{delayFlag},
				Action: func(ctx context.Context, c *cli.Command) error { return app.SelectClipboard(ctx, int(c.Int("delay"))) },
			},
			{
				Name: "select-file", Usage: "Capture selection to file",
				Flags:  []cli.Flag{delayFlag},
				Action: func(ctx context.Context, c *cli.Command) error { return app.SelectFile(ctx, int(c.Int("delay"))) },
			},
			{
				Name: "select-edit", Usage: "Capture selection and edit with satty",
				Flags:  []cli.Flag{delayFlag},
				Action: func(ctx context.Context, c *cli.Command) error { return app.SelectEdit(ctx, int(c.Int("delay"))) },
			},
			{
				Name: "window-clipboard", Usage: "Capture focused window to clipboard",
				Flags:  []cli.Flag{delayFlag},
				Action: func(ctx context.Context, c *cli.Command) error { return app.WindowClipboard(ctx, int(c.Int("delay"))) },
			},
			{
				Name: "window-file", Usage: "Capture focused window to file",
				Flags:  []cli.Flag{delayFlag},
				Action: func(ctx context.Context, c *cli.Command) error { return app.WindowFile(ctx, int(c.Int("delay"))) },
			},
			{
				Name: "screen-clipboard", Usage: "Capture focused screen to clipboard",
				Flags:  []cli.Flag{delayFlag},
				Action: func(ctx context.Context, c *cli.Command) error { return app.ScreenClipboard(ctx, int(c.Int("delay"))) },
			},
			{
				Name: "screen-file", Usage: "Capture focused screen to file",
				Flags:  []cli.Flag{delayFlag},
				Action: func(ctx context.Context, c *cli.Command) error { return app.ScreenFile(ctx, int(c.Int("delay"))) },
			},
			// Recording
			{
				Name: "record-select", Usage: "Record video of selection",
				Flags:  []cli.Flag{delayFlag},
				Action: func(ctx context.Context, c *cli.Command) error { return app.RecordSelect(ctx, int(c.Int("delay"))) },
			},
			{
				Name: "record-screen", Usage: "Record video of screen",
				Flags:  []cli.Flag{delayFlag},
				Action: func(ctx context.Context, c *cli.Command) error { return app.RecordScreen(ctx, int(c.Int("delay"))) },
			},
			{
				Name: "record-stop", Usage: "Stop recording and convert to mp4",
				Action: func(ctx context.Context, _ *cli.Command) error { return app.RecordStop(ctx) },
			},
			{
				Name: "record-pause", Usage: "Pause/resume recording",
				Action: func(ctx context.Context, _ *cli.Command) error { return app.RecordPause(ctx) },
			},
			{
				Name: "record-toggle", Usage: "Toggle recording (start select if idle, stop if recording)",
				Flags:  []cli.Flag{delayFlag},
				Action: func(ctx context.Context, c *cli.Command) error { return app.RecordToggle(ctx, int(c.Int("delay"))) },
			},
			// Waybar
			{
				Name: "waybar-status", Usage: "Output waybar status JSON",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "follow", Usage: "Continuously output on state change"},
				},
				Action: func(_ context.Context, c *cli.Command) error { return app.WaybarStatusCmd(c.Bool("follow")) },
			},
		},
	}

	if err := root.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
