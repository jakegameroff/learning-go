package main

import (
	"github.com/gorilla/websocket"
)

var hub *Hub

type Hub struct {
	clients map[*websocket.Conn]bool
	broadcast chan []byte
}

func initHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func run(hub *Hub) {
	for {
		msg := <- hub.broadcast // wait to receive a message from someone

		// when we do, broadcast it to everyone
		for conn := range hub.clients {
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func addClient(hub *Hub, conn *websocket.Conn) {
	hub.clients[conn] = true
}

func removeClient(hub *Hub, conn *websocket.Conn) {
	delete(hub.clients, conn)
	conn.Close()
}