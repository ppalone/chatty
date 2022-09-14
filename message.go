package main

type Message struct {
	Body string `json:"body"`
	By   string `json:"by"`
	Room string `json:"room"`
}
