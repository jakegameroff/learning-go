package main

import (
	"net/http"
	"github.com/gorilla/websocket"
)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{} // http --> ws protocol
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {return}

	addClient(hub, conn)
	defer removeClient(hub,conn)
	listenForClientMessage(conn)
}

func listenForClientMessage(conn *websocket.Conn) {
	for {
		// wait for a message
		_, msg, err := conn.ReadMessage()
		if err != nil {return}
		hub.broadcast <- msg // send message to the channel
	}
}