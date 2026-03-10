package hub

import "sync"

var (
	rooms   = make(map[string]*Hub)
	roomsMu sync.Mutex
)

func getRoom(roomName string) *Hub {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	hub, exists := rooms[roomName]
	if !exists {
		hub = initHub()
		rooms[roomName] = hub
		go run(hub)
	}
	return hub
}
