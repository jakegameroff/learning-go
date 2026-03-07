package main

var rooms = make(map[string]*Hub)

func getRoom(roomName string) *Hub {
	hub, exists := rooms[roomName]

	if !exists {
		hub = initHub()
		rooms[roomName] = hub
		go run(hub)
	}
	return hub
}
