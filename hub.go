package main

import (
	"log"
	"sync"
)

// Hub.
type Hub struct {
	// clients
	Clients map[*Client]bool

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
	h.Clients[c] = true
	h.mutex.Unlock()
	log.Println("Added client to Hub", h.Clients)
}

// Remove client from Hub.
func (h *Hub) delete(c *Client) {
	h.mutex.Lock()
	delete(h.Clients, c)
	h.mutex.Unlock()
	log.Println("Removed client from Hub")
}

// Broadcast message to all connected clients.
func (h *Hub) broadcast(m *Message) {
	for k := range h.Clients {
		k.Send <- m
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
