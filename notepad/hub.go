package main

import (
	"github.com/gorilla/websocket"
	"encoding/json"
)

type Hub struct {
	clients map[*websocket.Conn]string
	broadcast chan []byte
	history [][]byte
	cursorsChannel chan Cursor
	currentCursors map[string]int // userID --> position
	join chan ClientJoin
	leave chan *websocket.Conn
}

type Cursor struct {
	UserID string
	Position int
}

type ClientJoin struct {
	conn *websocket.Conn
	userID string
}

func initHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]string),
		broadcast: make(chan []byte),
		cursorsChannel: make(chan Cursor),
		currentCursors: make(map[string]int),
		join: make(chan ClientJoin),
		leave: make(chan *websocket.Conn),
	}
}

func run(hub *Hub) {
	for {
		select {
			case msg := <- hub.broadcast: // wait to receive a message from someone
				hub.history = append(hub.history, msg)

				// when we do, broadcast it to everyone
				for conn := range hub.clients {
					conn.WriteMessage(websocket.BinaryMessage, msg)
				}

			case cursor := <- hub.cursorsChannel:
				if cursor.Position == -1 {
					delete(hub.currentCursors, cursor.UserID)
				} else {
					hub.currentCursors[cursor.UserID] = cursor.Position
				}
				data, _ := json.Marshal(cursor)
				for conn := range hub.clients {
					if cursor.UserID != hub.clients[conn] {
						conn.WriteMessage(websocket.TextMessage, data)
					}
				}
			
			case client := <- hub.join:
				hub.clients[client.conn] = client.userID

				// send the history to the client
				for _, msg := range hub.history {
					client.conn.WriteMessage(websocket.BinaryMessage, msg)
				}
				// send cursors to the client
				for userID, pos := range hub.currentCursors {
					data, _ := json.Marshal(Cursor{UserID: userID, Position: pos})
					client.conn.WriteMessage(websocket.TextMessage, data)
				}

			case leave := <- hub.leave:
				userID, _ := hub.clients[leave]
				delete(hub.clients, leave)
				delete(hub.currentCursors, userID)
				data, _ := json.Marshal(Cursor{UserID: userID, Position: -1})

				// inform everyone else (cursor must be deleted!)
				for conn := range hub.clients {
					conn.WriteMessage(websocket.TextMessage, data)
				}
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
