package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type hub struct {
	// the mutex to protect connections
	connectionsMx sync.RWMutex

	// Registered connections.
	connections map[*connection]struct{}

	// Registered host connections
	hostConnections map[*connection]struct{}

	// admin messages
	adminChan chan []byte

	// Inbound messages from the connections.
	broadcast chan []byte

	logMx sync.RWMutex
	log   [][]byte
}

func newHub() *hub {
	h := &hub{
		connectionsMx:   sync.RWMutex{},
		broadcast:       make(chan []byte),
		adminChan:       make(chan []byte),
		connections:     make(map[*connection]struct{}),
		hostConnections: make(map[*connection]struct{}),
	}

	go func() { // this runs only on app startup
		for { // this runs when something is added to broadcast channel
			// 3. this is third
			msg := <-h.broadcast // 1. this is first - listening on the broadcast channel
			h.connectionsMx.RLock()
			for c := range h.connections { // 2. this is second
				select {
				case c.send <- msg: // send to each of the connections send channels, which the writer is listening on
				// stop trying to send to this connection after trying for 1 second.
				// if we have to stop, it means that a reader died so remove the connection also.
				case <-time.After(1 * time.Second):
					log.Printf("shutting down connection %s", c)
					h.removeConnection(c)
				}
			}
			h.connectionsMx.RUnlock()
		}

	}()

	go func() {
		for { // this runs when something is added to broadcast channel
			// 3. this is third
			msg := <-h.adminChan // 1. this is first - listening on the broadcast channel
			h.connectionsMx.RLock()
			for c := range h.hostConnections { // 2. this is second
				select {
				case c.send <- msg: // send to each of the connections send channels, which the writer is listening on
				// stop trying to send to this connection after trying for 1 second.
				// if we have to stop, it means that a reader died so remove the connection also.
				case <-time.After(1 * time.Second):
					log.Printf("shutting down connection %s", c)
					h.removeConnection(c)
				}
			}
			h.connectionsMx.RUnlock()
		}
	}()
	return h
}

func (h *hub) addConnection(conn *connection) {
	fmt.Println("adding connection")
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	h.connections[conn] = struct{}{}
}

func (h *hub) removeConnection(conn *connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	if _, ok := h.connections[conn]; ok {
		delete(h.connections, conn)
		close(conn.send)
	}
}

func (h *hub) addHostConnection(conn *connection) {
	fmt.Println("adding host connection")
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	h.hostConnections[conn] = struct{}{}
}

func (h *hub) removeHostConnection(conn *connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	if _, ok := h.hostConnections[conn]; ok {
		delete(h.hostConnections, conn)
		close(conn.send)
	}
}
