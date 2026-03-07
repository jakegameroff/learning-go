package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Hub struct {
	game       Game
	players    map[*websocket.Conn]string
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func initHub() *Hub {
	return &Hub{
		game:       initGame(),
		players:    make(map[*websocket.Conn]string),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func run(h *Hub) {
	for {
		select {
		case conn := <-h.register:
			h.players[conn] = ""
		case conn := <-h.unregister:
			delete(h.players, conn)
			conn.Close()
		case hexMove := <-h.broadcast:
			var node Node
			json.Unmarshal(hexMove, &node)
			node.NodeType = "board_node"
			h.game.move(node)

			for conn := range h.players {
				conn.WriteMessage(websocket.TextMessage, hexMove)
			}

			if h.game.isWinningMove(node) {
				msg, _ := json.Marshal(map[string]string{"winner": node.Color})
				for conn := range h.players {
					conn.WriteMessage(websocket.TextMessage, msg)
				}
			}
		}
	}
}

func addClient(hub *Hub, conn *websocket.Conn) {
	hub.register <- conn
}

func removeClient(hub *Hub, conn *websocket.Conn) {
	hub.unregister <- conn
}
