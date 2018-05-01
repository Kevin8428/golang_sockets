package main

import (
	"fmt"
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

func (c *connection) reader(wg *sync.WaitGroup, wsConn *websocket.Conn, isAdmin bool) {
	defer wg.Done()
	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println("reading!")
		if isAdmin {
			c.h.adminChan <- message
		} else {
			c.h.broadcast <- message
		}
	}
}

func loadPreviousMessages(wsConn *websocket.Conn) {
	for _, message := range previousMessages {
		_ = wsConn.WriteMessage(websocket.TextMessage, message)
	}
}

func (c *connection) writer(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for message := range c.send { // listening to send channel
		fmt.Println("one message being sent")
		err := wsConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type customerWSHandler struct {
	h *hub
}

func (wsh customerWSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading %s", err)
		return
	}
	c := &connection{send: make(chan []byte, 256), h: wsh.h}
	c.h.addConnection(c)
	defer c.h.removeConnection(c)
	var wg sync.WaitGroup
	wg.Add(2)
	loadPreviousMessages(wsConn)
	go c.writer(&wg, wsConn)
	go c.reader(&wg, wsConn, false)
	wg.Wait()
	wsConn.Close()
}

type adminWSHandler struct {
	h *hub
}

func (wsh adminWSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading %s", err)
		return
	}
	c := &connection{send: make(chan []byte, 256), h: wsh.h}
	c.h.addHostConnection(c)
	defer c.h.removeHostConnection(c)
	var wg sync.WaitGroup
	wg.Add(2)
	go c.writer(&wg, wsConn)
	go c.reader(&wg, wsConn, true)
	wg.Wait()
	wsConn.Close()
}
