package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type connection struct {
	// Buffered channel of outbound messages.
	send chan []byte

	// The hub.
	h *hub
}

var previousMessages = [][]byte{
	[]byte{104, 105},
	[]byte{104, 101, 121},
	[]byte{121, 111},
}

func (c *connection) reader(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
		c.h.broadcast <- message
	}
}

func loadPreviousMessages(wsConn *websocket.Conn) {
	_ = wsConn.WriteMessage(websocket.TextMessage, previousMessages[0])
	_ = wsConn.WriteMessage(websocket.TextMessage, previousMessages[1])
	_ = wsConn.WriteMessage(websocket.TextMessage, previousMessages[2])
}

func (c *connection) writer(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for message := range c.send { // listening to send channel
		err := wsConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type wsHandler struct {
	h *hub
}

func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading %s", err)
		return
	}
	c := &connection{send: make(chan []byte, 256), h: wsh.h}
	c.h.addConnection(c)
	defer c.h.removeConnection(c)
	var wg sync.WaitGroup
	wg.Add(3)
	loadPreviousMessages(wsConn)
	go c.writer(&wg, wsConn)
	go c.reader(&wg, wsConn)
	wg.Wait()
	wsConn.Close()
}