package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader *websocket.Upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ws(w http.ResponseWriter, r *http.Request, hub *Hub, room string) {
	// upgrade http request to websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading request to websocket connection: ", err)
	}

	// new client
	client := &Client{
		Conn: conn,
		Send: make(chan *Message),
		Hub:  hub,
		Room: room,
	}

	// register client to hub
	client.Hub.Register <- client

	go client.Read()
	go client.Write()
}

func main() {
	port := os.Getenv("PORT")
	if strings.Trim(port, " ") == "" {
		port = "8000"
	}
	hub := &Hub{
		Clients:    make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
		mutex:      &sync.RWMutex{},
	}
	go hub.Run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/general", http.StatusTemporaryRedirect)
			return
		}
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		t := strings.Split(strings.Trim(r.URL.Path, " "), "/")
		room := "general"
		if len(t) > 2 && t[2] != "" {
			room = t[2]
		}
		ws(w, r, hub, room)
	})
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
