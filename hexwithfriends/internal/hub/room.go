package hub

import (
	"strings"
	"sync"
)

var (
	rooms   = make(map[string]*Hub)
	roomsMu sync.Mutex
)

func getRoom(roomName string) *Hub {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	roomName = strings.ToLower(roomName)
	hub, exists := rooms[roomName]
	if !exists {
		hub = initHub()
		rooms[roomName] = hub
		go run(hub)
	}
	return hub
}
