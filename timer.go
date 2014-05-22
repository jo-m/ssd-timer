package main

import (
	"time"
)

// holds the static config for a timer
type TimerConfig struct {
	NumRounds      int
	PauseLength    time.Duration
	RoundLength    time.Duration
	Title          string
	RoundText      string
	PauseText      string
	FinishedText   string
	NotStartedText string
}

// holds the current global state for the timer.
// These values are only changed when pressing pause, run etc.
type TimerState struct {
	TimerStartAt  time.Time
	TimerPaused   bool
	TimerPausedAt time.Time
}

type Timer struct {
	Config TimerConfig
	State  TimerState
}
