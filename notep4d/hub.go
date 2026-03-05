package main

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients   map[*websocket.Conn]string
	broadcast chan []byte
	history   [][]byte
	relay     chan RelayMsg
	join      chan ClientJoin
	leave     chan *websocket.Conn
}

type RelayMsg struct {
	sender *websocket.Conn
	data   []byte
}

type ClientJoin struct {
	conn   *websocket.Conn
	userID string
}

func initHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]string),
		broadcast: make(chan []byte),
		relay:     make(chan RelayMsg),
		join:      make(chan ClientJoin),
		leave:     make(chan *websocket.Conn),
	}
}

func run(hub *Hub) {
	for {
		select {
		case msg := <-hub.broadcast:
			hub.history = append(hub.history, msg)
			for conn := range hub.clients {
				conn.WriteMessage(websocket.BinaryMessage, msg)
			}

		case relay := <-hub.relay:
			for conn := range hub.clients {
				if conn != relay.sender {
					conn.WriteMessage(websocket.TextMessage, relay.data)
				}
			}

		case client := <-hub.join:
			hub.clients[client.conn] = client.userID
			for _, msg := range hub.history {
				client.conn.WriteMessage(websocket.BinaryMessage, msg)
			}
			// indicate to client that we are done sending the history 
			client.conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"ready"}`))

		case leave := <-hub.leave:
			delete(hub.clients, leave)
			leave.Close()
		}
	}
}

func addClient(hub *Hub, conn *websocket.Conn, userID string) {
	hub.join <- ClientJoin{conn: conn, userID: userID}
}

func removeClient(hub *Hub, conn *websocket.Conn) {
	hub.leave <- conn
}
