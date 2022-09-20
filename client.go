package main

import (
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

// Client.
type Client struct {
	// websocket connection
	Conn *websocket.Conn

	// send channel
	Send chan *WSMessage

	// Hub
	Hub *Hub

	// Room name
	Room string

	// username
	Username string
}

type WSMessage struct {
	Type    string  `json:"type"`
	Payload Message `json:"payload"`
}

func (c *Client) Read() {
	defer func() {
		// unregister client
		c.Hub.Unregister <- c
		// close channel
		close(c.Send)
	}()
	for {
		var message WSMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			log.Println("Error while reading websocket message: ", err)
			return
		}
		log.Println(message)
		switch message.Type {

		case "join":
			log.Println("User joined", message.Payload)
			c.Username = message.Payload.By
			var m *WSMessage = &WSMessage{
				Type: "join",
				Payload: Message{
					Body: strconv.Itoa(len(c.Hub.Clients[message.Payload.Room])),
					By:   message.Payload.By,
					Room: message.Payload.Room,
				},
			}
			c.Hub.Broadcast <- m

		case "message":
			var m *WSMessage = &WSMessage{
				Type: "message",
				Payload: Message{
					Body: message.Payload.Body,
					By:   message.Payload.By,
					Room: message.Payload.Room,
				},
			}
			c.Hub.Broadcast <- m

		case "typing":
			var m *WSMessage = &WSMessage{
				Type: message.Type,
				Payload: Message{
					By:   c.Username,
					Room: message.Payload.Room,
				},
			}
			log.Println("Typing:", m)
			c.Hub.Broadcast <- m

		case "stoptyping":
			var m *WSMessage = &WSMessage{
				Type: message.Type,
				Payload: Message{
					By:   c.Username,
					Room: message.Payload.Room,
				},
			}
			log.Println("StopTyping:", m)
			c.Hub.Broadcast <- m
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		m, ok := <-c.Send
		if !ok {
			log.Println("Send Channel was closed")
			return
		}
		log.Println("Write message:", m)
		// var message *WSMessage = &WSMessage{
		// 	Type:    "message",
		// 	Payload: *m,
		// }
		err := c.Conn.WriteJSON(m)
		if err != nil {
			log.Println("Error while writing message:", err)
			return
		}
	}
}
