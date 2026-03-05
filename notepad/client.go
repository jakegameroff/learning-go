package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"math/rand"
	"encoding/json"
	"os"
)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Path
	hub := getRoom(roomName)

	var upgrader = websocket.Upgrader{} // http --> ws protocol
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {return}

	// make user id
	data, _ := os.ReadFile("animals.json")
	var animals []string
	json.Unmarshal(data, &animals)
	userID := "Anonymous " + animals[rand.Intn(len(animals))]

	addClient(hub, conn, userID)
	defer removeClient(hub,conn)
	listenForClientMessage(hub, conn, userID)
}

func listenForClientMessage(hub *Hub, conn *websocket.Conn, userID string) {
	for {
		// wait for a message
		msgType, msg, err := conn.ReadMessage()
		if err != nil {return}

		if msgType == websocket.TextMessage {
			var cursor Cursor
			json.Unmarshal(msg, &cursor)
			cursor.UserID = userID
			hub.cursorsChannel <- cursor

		} else {
			hub.broadcast <- msg // send message to the channel
		}
	}
}
