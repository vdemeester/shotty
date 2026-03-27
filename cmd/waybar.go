package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vdemeester/shotty/internal/state"
)

// WaybarStatus is the JSON structure expected by waybar's custom module.
type WaybarStatus struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
	Alt     string `json:"alt"`
}

func buildWaybarStatus(s *state.State) *WaybarStatus {
	if s.Countdown > 0 {
		return &WaybarStatus{
			Text:    fmt.Sprintf("⏱ %d", s.Countdown),
			Tooltip: fmt.Sprintf("Starting in %ds", s.Countdown),
			Class:   "countdown",
			Alt:     "countdown",
		}
	}

	if s.Recording {
		elapsed := time.Since(s.StartedAt)
		minutes := int(elapsed.Minutes())
		seconds := int(elapsed.Seconds()) % 60
		timeStr := fmt.Sprintf("%02d:%02d", minutes, seconds)

		if s.Paused {
			return &WaybarStatus{
				Text:    fmt.Sprintf("⏸ %s", timeStr),
				Tooltip: "Recording paused",
				Class:   "paused",
				Alt:     "paused",
			}
		}

		return &WaybarStatus{
			Text:    fmt.Sprintf("⏺ %s", timeStr),
			Tooltip: fmt.Sprintf("Recording: %s", s.File),
			Class:   "recording",
			Alt:     "recording",
		}
	}

	return &WaybarStatus{
		Text:    "",
		Tooltip: "",
		Class:   "idle",
		Alt:     "idle",
	}
}

// BuildCurrentStatus reads state and returns the current waybar status.
func (a *App) BuildCurrentStatus() (*WaybarStatus, error) {
	s, err := state.Read(a.Config.StateFile)
	if err != nil {
		s = &state.State{}
	}
	return buildWaybarStatus(s), nil
}

// WaybarStatus outputs waybar-compatible JSON. With follow=true, it polls
// every second and prints on change.
func (a *App) WaybarStatusCmd(follow bool) error {
	if !follow {
		status, _ := a.BuildCurrentStatus()
		return json.NewEncoder(os.Stdout).Encode(status)
	}

	// Follow mode: poll every second, print on change
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var previous string

	for {
		select {
		case <-ticker.C:
			s, err := state.Read(a.Config.StateFile)
			if err != nil {
				s = &state.State{}
			}
			status := buildWaybarStatus(s)
			data, _ := json.Marshal(status)
			current := string(data)

			// Always print during recording (elapsed time changes)
			if current != previous || s.Recording {
				fmt.Println(current)
				previous = current
			}
		case <-sigChan:
			return nil
		}
	}
}
