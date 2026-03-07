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
	playerID := r.URL.Query().Get("name")
	hub := getRoom(roomName)

	var upgrader = websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	if playerID == "" {
		var names []string
		data, _ := os.ReadFile("static/names.json")
		json.Unmarshal(data, &names)
		playerID = names[rand.Intn(len(names))]
	}

	idMsg, _ := json.Marshal(map[string]string{"playerID": playerID})
	conn.WriteMessage(websocket.TextMessage, idMsg)

	addClient(hub, conn)
	defer removeClient(hub, conn)
	listenForClientMessage(hub, conn)
}

func listenForClientMessage(hub *Hub, conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		hub.broadcast <- msg
	}
}
