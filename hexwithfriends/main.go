package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		upgradeHeader := r.Header.Get("Upgrade")
		if upgradeHeader == "websocket" {
			handleWebSocket(w, r)
		} else if r.URL.Path == "/" {
			http.ServeFile(w, r, "static/landing.html")
		} else {
			http.ServeFile(w, r, "static/index.html")
		}
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
