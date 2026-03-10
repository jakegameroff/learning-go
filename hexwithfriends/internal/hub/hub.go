package hub

import (
	"encoding/json"
	"hexwithfriends/internal/game"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	conn        *websocket.Conn
	playerName  string
	sessionID   string
	isConnected bool

	color  string
	myTurn bool
}

type HexMove struct {
	conn *websocket.Conn
	data []byte
}

type Hub struct {
	game          game.Game
	players       map[string]User
	connToSession map[*websocket.Conn]string
	broadcast     chan HexMove
	register      chan User
	unregister    chan User
}

func initHub() *Hub {
	return &Hub{
		game:          game.InitGame(),
		players:       make(map[string]User),
		connToSession: make(map[*websocket.Conn]string),
		broadcast:     make(chan HexMove),
		register:      make(chan User),
		unregister:    make(chan User),
	}
}

func run(h *Hub) {
	for {
		select {
		case user := <-h.register:
			h.handleRegister(user)
		case user := <-h.unregister:
			h.handleUnregister(user)
		case hexMove := <-h.broadcast:
			h.handleMove(hexMove)
		}
	}
}

func (h *Hub) handleReconnect(user User, existing User) {
	existing.conn = user.conn
	existing.isConnected = true
	h.players[user.sessionID] = existing
	h.connToSession[user.conn] = user.sessionID

	m, _ := json.Marshal(map[string]interface{}{
		"reconnected": true, "color": existing.color,
		"playerName": existing.playerName, "sessionID": user.sessionID,
		"board":     h.game.Board[:game.PlayableCells],
		"isRedTurn": h.game.IsRedTurn,
	})
	user.conn.WriteMessage(websocket.TextMessage, m)

	h.broadcastMessage(map[string]string{"opponent_reconnected": "true"})
}

func (h *Hub) handleRegister(user User) {
	if existing, ok := h.players[user.sessionID]; ok && !existing.isConnected {
		h.handleReconnect(user, existing)
		return
	}

	if len(h.players) >= 2 {
		// TODO: add spectator mode!
		msg, _ := json.Marshal(map[string]string{"error": "room_full"})
		user.conn.WriteMessage(websocket.TextMessage, msg)
		user.conn.Close()
		return
	}

	user.sessionID = uuid.New().String()
	user.isConnected = true
	assignColor(&user, h)
	h.connToSession[user.conn] = user.sessionID

	msg, _ := json.Marshal(map[string]string{"color": user.color})
	user.conn.WriteMessage(websocket.TextMessage, msg)

	if len(h.players) == 2 {
		for sid, u := range h.players {
			m, _ := json.Marshal(map[string]interface{}{
				"start": "true", "color": u.color,
				"playerName": u.playerName, "sessionID": u.sessionID,
			})
			if u.conn != nil {
				u.conn.WriteMessage(websocket.TextMessage, m)
			}
			_ = sid
		}
	}
}

func (h *Hub) handleUnregister(user User) {
	sid, ok := h.connToSession[user.conn]
	if !ok {
		user.conn.Close()
		return
	}
	delete(h.connToSession, user.conn)
	u := h.players[sid]
	u.conn = nil
	u.isConnected = false
	h.players[sid] = u
	user.conn.Close()
}

func (h *Hub) handleMove(hexMove HexMove) {
	sid, ok := h.connToSession[hexMove.conn]
	if !ok {
		return
	}
	player := h.players[sid]

	if !player.myTurn {
		msg, _ := json.Marshal(map[string]string{"error": "not_your_turn"})
		hexMove.conn.WriteMessage(websocket.TextMessage, msg)
		return
	}

	var node game.Node
	json.Unmarshal(hexMove.data, &node)
	node.IsWinNode = false
	node.Color = player.color
	result := h.game.Move(node)

	if !result {
		msg, _ := json.Marshal(map[string]string{"error": "invalid_move"})
		hexMove.conn.WriteMessage(websocket.TextMessage, msg)
		return
	}

	h.swapTurns()
	h.broadcastMessage(map[string]interface{}{"index": node.Index, "color": node.Color})

	if h.game.IsWinningMove(node) {
		h.handleWin(node)
	}
}

func (h *Hub) swapTurns() {
	for sid, u := range h.players {
		u.myTurn = !u.myTurn
		h.players[sid] = u
	}
}

func (h *Hub) handleWin(node game.Node) {
	h.broadcastMessage(map[string]string{"winner": node.Color})

	for sid, u := range h.players {
		if u.color == "red" {
			u.color = "blue"
			u.myTurn = false
		} else {
			u.color = "red"
			u.myTurn = true
		}
		h.players[sid] = u
	}

	h.game = game.InitGame()

	for _, u := range h.players {
		m, _ := json.Marshal(map[string]interface{}{"reset": true, "color": u.color})
		if u.conn != nil {
			u.conn.WriteMessage(websocket.TextMessage, m)
		}
	}
}

func (h *Hub) getBoardState() []game.Node {
	return h.game.Board[:]
}

func (h *Hub) broadcastMessage(msg interface{}) {
	data, _ := json.Marshal(msg)
	for _, u := range h.players {
		if u.conn != nil && u.isConnected {
			u.conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func addClient(hub *Hub, conn *websocket.Conn, playerName string) {
	hub.register <- User{conn: conn, playerName: playerName}
}

func removeClient(hub *Hub, conn *websocket.Conn) {
	hub.unregister <- User{conn: conn}
}
