package hub

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Path
	playerName := r.URL.Query().Get("name")
	sessionID := r.URL.Query().Get("session")
	hub := getRoom(roomName)

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	if sessionID != "" {
		hub.register <- User{conn: conn, sessionID: sessionID}
		defer removeClient(hub, conn)
		listenForClientMessage(hub, conn)
		return
	}

	if playerName == "" {
		var names []string
		data, _ := os.ReadFile("static/names.json")
		json.Unmarshal(data, &names)
		playerName = names[rand.Intn(len(names))]
	}

	idMsg, _ := json.Marshal(map[string]string{"playerName": playerName})
	conn.WriteMessage(websocket.TextMessage, idMsg)

	addClient(hub, conn, playerName)
	defer removeClient(hub, conn)
	listenForClientMessage(hub, conn)
}

func listenForClientMessage(hub *Hub, conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		hub.broadcast <- HexMove{conn: conn, data: msg}
	}
}
