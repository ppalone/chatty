package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader *websocket.Upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ws(w http.ResponseWriter, r *http.Request, hub *Hub) {
	// upgrade http request to websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading request to websocket connection: ", err)
	}

	// new client
	client := &Client{
		Conn: conn,
		Send: make(chan Message),
		Hub:  hub,
	}

	// register client to hub
	client.Hub.Register <- *client

	go client.Read()
}

func main() {
	hub := &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan Client),
		Unregister: make(chan Client),
		Broadcast:  make(chan Message),
		mutex:      &sync.RWMutex{},
	}
	go hub.Run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws(w, r, hub)
	})
	log.Fatalln(http.ListenAndServe(":8000", nil))
}
