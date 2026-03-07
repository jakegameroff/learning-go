package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type User struct {
	conn     *websocket.Conn
	playerID string
	color    string
	myTurn   bool
}

type HexMove struct {
	conn *websocket.Conn
	data []byte
}

type Hub struct {
	game       Game
	players    map[*websocket.Conn]User
	broadcast  chan HexMove
	register   chan User
	unregister chan User
}

func initHub() *Hub {
	return &Hub{
		game:       initGame(),
		players:    make(map[*websocket.Conn]User),
		broadcast:  make(chan HexMove),
		register:   make(chan User),
		unregister: make(chan User),
	}
}

// TODO: validate moves — check index is in bounds and cell is unoccupied

func run(h *Hub) {
	for {
		select {
		case user := <-h.register:
			if len(h.players) >= 2 {
				msg, _ := json.Marshal(map[string]string{"error": "room_full"})
				user.conn.WriteMessage(websocket.TextMessage, msg)
				user.conn.Close()
				continue
			}

			assignColor(&user, h)
			msg, _ := json.Marshal(map[string]string{"color": user.color})
			user.conn.WriteMessage(websocket.TextMessage, msg)

			if len(h.players) == 2 {
				start, _ := json.Marshal(map[string]string{"start": "true"})
				for conn := range h.players {
					conn.WriteMessage(websocket.TextMessage, start)
				}
			}

		case user := <-h.unregister:
			delete(h.players, user.conn)
			user.conn.Close()

		case hexMove := <-h.broadcast:
			player := h.players[hexMove.conn]

			if !player.myTurn {
				msg, _ := json.Marshal(map[string]string{"error": "not_your_turn"})
				hexMove.conn.WriteMessage(websocket.TextMessage, msg)
				continue
			}

			var node Node
			json.Unmarshal(hexMove.data, &node)
			node.IsWinNode = false
			node.Color = player.color
			result := h.game.move(node)

			if !result {
				msg, _ := json.Marshal(map[string]string{"error": "invalid_move"})
				hexMove.conn.WriteMessage(websocket.TextMessage, msg)
				continue
			}

			for conn, u := range h.players {
				u.myTurn = !u.myTurn
				h.players[conn] = u
			}

			moveMsg, _ := json.Marshal(
				map[string]interface{}{
					"index": node.Index, "color": node.Color,
				})
			for conn := range h.players {
				conn.WriteMessage(websocket.TextMessage, moveMsg)
			}

			if h.game.isWinningMove(node) {
				msg, _ := json.Marshal(map[string]string{"winner": node.Color})
				for conn := range h.players {
					conn.WriteMessage(websocket.TextMessage, msg)
					u := h.players[conn]
					if u.color == "red" {
						u.color = "blue"
						u.myTurn = false
					} else {
						u.color = "red"
						u.myTurn = true
					}
					h.players[conn] = u
				}
				h.game = initGame()
				for conn, u := range h.players {
					m, _ := json.Marshal(map[string]interface{}{"reset": true, "color": u.color})
					conn.WriteMessage(websocket.TextMessage, m)
				}
			}
		}
	}
}

func addClient(hub *Hub, conn *websocket.Conn, playerID string) {
	hub.register <- User{conn: conn, playerID: playerID}
}

func removeClient(hub *Hub, conn *websocket.Conn) {
	hub.unregister <- User{conn: conn}
}
