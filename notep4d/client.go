package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Path
	hub := getRoom(roomName)

	var upgrader = websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// make user id
	data, _ := os.ReadFile("static/animals.json")
	var animals []string
	json.Unmarshal(data, &animals)
	userID := "Anonymous " + animals[rand.Intn(len(animals))]

	// send userID to client
	idMsg, _ := json.Marshal(map[string]string{"userID": userID})
	conn.WriteMessage(websocket.TextMessage, idMsg)

	addClient(hub, conn, userID)
	defer removeClient(hub, conn)
	listenForClientMessage(hub, conn)
}

func listenForClientMessage(hub *Hub, conn *websocket.Conn) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if msgType == websocket.TextMessage {
			hub.relay <- RelayMsg{sender: conn, data: msg}
		} else {
			hub.broadcast <- msg
		}
	}
}
