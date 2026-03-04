package main

import (
	"net/http"
)

func main() {
	hub = initHub()
	go run(hub)

	http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8000", nil)
}