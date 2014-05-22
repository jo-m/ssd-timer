package main

import (
	"github.com/gorilla/websocket"
	"time"
)

const (
	// Time allowed to write a message
	writeWait = 10 * time.Second

	// Time allowed to read pong answer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 8
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type connection struct {
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Timer
}

// write writes a message with the given message type and payload
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}
