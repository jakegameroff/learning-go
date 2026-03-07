package main

import (
	"math/rand"
)

func assignColor(user *User, h *Hub) {
	if len(h.players) == 0 {
		if rand.Intn(2) == 0 {
			user.color = "red"
			user.myTurn = true
		} else {
			user.color = "blue"
			user.myTurn = false
		}
		h.players[user.conn] = *user
	} else if len(h.players) == 1 {
		var u User
		for _, u = range h.players {
			break
		}
		switch u.color {
		case "red":
			user.color = "blue"
			user.myTurn = false
		default:
			user.color = "red"
			user.myTurn = true
		}
		h.players[user.conn] = *user
	} else {
		return
	}
}

func (g *Game) move(node Node) bool {
	if node.Index < 0 || node.Index >= size*size || g.Board[node.Index].Color != "" {
		return false
	}
	g.Board[node.Index] = node
	return true
}
