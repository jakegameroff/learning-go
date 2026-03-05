package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		upgradeHeader := r.Header.Get("Upgrade")
		if upgradeHeader == "websocket" {
			handleWebSocket(w, r)
		} else {
			http.ServeFile(w, r, "index.html")
		}
	})

	http.ListenAndServe(":8000", nil)
}