package main

import (
	"encoding/json"
	"log"
	"time"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	clientConnections map[*connection]bool
	adminConnections  map[*connection]bool

	// chan the admin ws sends its message to
	// so it gets parsed and put to timer
	adminCommands chan []byte

	timer Timer

	registerClient   chan *connection
	unregisterClient chan *connection
	registerAdmin    chan *connection
	unregisterAdmin  chan *connection
}

var defaultTimer = Timer{
	Config: TimerConfig{
		NumRounds:      14,
		PauseLength:    time.Second * 30,
		RoundLength:    time.Minute * 5,
		Title:          "SSD Timer",
		RoundText:      "ends in",
		PauseText:      "starts in",
		FinishedText:   "finished",
		NotStartedText: "not started yet",
	},
	State: TimerState{
		TimerStartAt:  time.Now(),
		TimerPaused:   true,
		TimerPausedAt: time.Now(),
	},
}

var h = hub{
	clientConnections: make(map[*connection]bool),
	adminConnections:  make(map[*connection]bool),
	adminCommands:     make(chan []byte, 1024),

	timer: defaultTimer,

	registerClient:   make(chan *connection),
	unregisterClient: make(chan *connection),
	registerAdmin:    make(chan *connection),
	unregisterAdmin:  make(chan *connection),
}

func (h *hub) sendState(c *connection) {
	select {
	case c.send <- h.timer:
	default:
		close(c.send)
		delete(h.adminConnections, c)
		delete(h.clientConnections, c)
	}
}

func (h *hub) broadcastStateAndConfig() {
	for c := range h.clientConnections {
		h.sendState(c)
	}
	for c := range h.adminConnections {
		h.sendState(c)
	}
}

func (h *hub) executeAdminCommand(cmd string) {
	log.Println("admin cmd:", cmd)
	switch cmd {
	case "round--":
		if !h.timer.State.TimerPaused {
			h.timer.State.TimerStartAt = h.timer.State.TimerStartAt.Add(h.timer.Config.PauseLength + h.timer.Config.RoundLength)
		}
	case "round++":
		if !h.timer.State.TimerPaused {
			h.timer.State.TimerStartAt = h.timer.State.TimerStartAt.Add(-(h.timer.Config.PauseLength + h.timer.Config.RoundLength))
		}
	case "minute--":
		if !h.timer.State.TimerPaused {
			h.timer.State.TimerStartAt = h.timer.State.TimerStartAt.Add(time.Minute)
		}
	case "minute++":
		if !h.timer.State.TimerPaused {
			h.timer.State.TimerStartAt = h.timer.State.TimerStartAt.Add(-time.Minute)
		}
	case "pause":
		if !h.timer.State.TimerPaused {
			h.timer.State.TimerPaused = true
			h.timer.State.TimerPausedAt = time.Now()
		}
	case "run":
		if h.timer.State.TimerPaused {
			h.timer.State.TimerPaused = false
			h.timer.State.TimerStartAt = h.timer.State.TimerStartAt.Add(time.Since(h.timer.State.TimerPausedAt))
		}
	case "reset":
		h.timer = defaultTimer
		h.timer.State.TimerStartAt = time.Now()
		h.timer.State.TimerPausedAt = time.Now()
	case "resetconfig":
		h.timer.Config = defaultTimer.Config
	default:
		newConfig := TimerConfig{}
		err := json.Unmarshal([]byte(cmd), &newConfig)
		if err == nil {
			h.timer.Config = newConfig
		}
	}
	h.broadcastStateAndConfig()
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.registerClient:
			h.clientConnections[c] = true

			h.sendState(c)
		case c := <-h.unregisterClient:
			delete(h.clientConnections, c)
			close(c.send)
		case c := <-h.registerAdmin:
			h.adminConnections[c] = true

			h.sendState(c)
		case c := <-h.unregisterAdmin:
			delete(h.adminConnections, c)
			close(c.send)
		case message := <-h.adminCommands:
			h.executeAdminCommand(string(message))
		}
	}
}
