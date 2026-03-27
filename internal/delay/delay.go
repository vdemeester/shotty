// Package delay implements a countdown timer that writes state for waybar display.
package delay

import (
	"time"

	"github.com/vdemeester/shotty/internal/state"
)

// Countdown sleeps for the given number of seconds, updating the state file
// each second so waybar can display the countdown.
func Countdown(stateFile string, seconds int) {
	if seconds <= 0 {
		return
	}
	for i := seconds; i > 0; i-- {
		_ = state.Write(stateFile, &state.State{Countdown: i})
		time.Sleep(time.Second)
	}
	_ = state.Clear(stateFile)
}
