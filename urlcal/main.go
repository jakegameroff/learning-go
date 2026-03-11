package main

import (
	"fmt"
	"net/http"
)

func main() {
	db := initDB()
	fmt.Println("Database initialized")

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		handleAPI(db, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "static/index.html")
			return
		}
		http.ServeFile(w, r, "static/calendar.html")
	})

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
