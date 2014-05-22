package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func (c *connection) adminReadLoop() {
	defer func() {
		h.unregisterAdmin <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		h.adminCommands <- message
	}
}

func (c *connection) adminWriteLoop() {
	log.Println("admin connect: ", c.ws.RemoteAddr())
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("admin disconnect: ", c.ws.RemoteAddr())
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

func serveAdminWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	cookie, _ := r.Cookie(cookieName)
	if cookie == nil || cookie.Value != *password {
		http.Error(w, "Login first!", 401)
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
	h.registerAdmin <- c
	go c.adminWriteLoop()
	c.adminReadLoop()
}
