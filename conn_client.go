package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func (c *connection) clientReadLoop() {
	defer func() {
		h.unregisterClient <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *connection) clientWriteLoop() {
	log.Println("client connect: ", c.ws.RemoteAddr())
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("client disconnect: ", c.ws.RemoteAddr())
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveTimerWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	c := &connection{send: make(chan Timer, 256), ws: ws}
	h.registerClient <- c
	go c.clientWriteLoop()
	c.clientReadLoop()
}
