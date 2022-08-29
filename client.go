package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client.
type Client struct {
	// websocket connection
	Conn *websocket.Conn

	// send channel
	Send chan Message

	// Hub
	Hub *Hub
}

type WSMessage struct {
	Type    string
	Payload interface{}
}

func (c *Client) Read() {
	defer func() {
		// unregister client
		c.Hub.Unregister <- *c
		// close channel
		close(c.Send)
		// close websocket connection
	}()
	for {
		var message WSMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			log.Println("Error while reading websocket message: ", err)
			break
		}
		switch message.Type {
		case "message":
			log.Println(message.Payload)
		}
	}
}
