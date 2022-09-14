package main

import (
	"log"
	"sync"
)

// Hub.
type Hub struct {
	// clients
	Clients map[string]map[*Client]bool

	// register channel
	Register chan *Client

	// unregister channel
	Unregister chan *Client

	// broadcast channel
	Broadcast chan *Message

	// mutex
	mutex *sync.RWMutex
}

// Add client to Hub.
func (h *Hub) add(c *Client) {
	h.mutex.Lock()
	if _, ok := h.Clients[c.Room]; !ok {
		h.Clients[c.Room] = make(map[*Client]bool)
	}
	h.Clients[c.Room][c] = true
	h.mutex.Unlock()
	log.Printf("Client added to room #%s, Number of clients in room: %d\n", c.Room, len(h.Clients[c.Room]))
}

// Remove client from Hub.
func (h *Hub) delete(c *Client) {
	h.mutex.Lock()
	if clients, ok := h.Clients[c.Room]; ok {
		delete(clients, c)
		c.Conn.Close()
	}
	h.mutex.Unlock()
	log.Printf("Removed client from room #%s, Number of clients in Room: %d\n", c.Room, len(h.Clients[c.Room]))
}

// Broadcast message to all connected clients.
func (h *Hub) broadcast(m *Message) {
	if clients, ok := h.Clients[m.Room]; ok {
		for k := range clients {
			k.Send <- m
		}
	}
}

// Run
func (h *Hub) Run() {
	log.Println("Hub is running")
	for {
		select {
		case client := <-h.Register:
			h.add(client)
		case client := <-h.Unregister:
			h.delete(client)
		case m := <-h.Broadcast:
			log.Println("Reached broadcast channel", m)
			h.broadcast(m)
		}
	}
}
