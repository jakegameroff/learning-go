package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Calendar struct {
	CalName  string     `json:"calName"`
	Entries  []CalEntry `json:"entries"`
	Password string     `json:"password"`
}

type CalEntry struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Start string `json:"start"`
	End   string `json:"end"`
	Note  string `json:"note"`
	Color string `json:"color"`
}

func getAuth(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if len(token) > len("Bearer ") {
		return token[len("Bearer "):]
	}
	return ""
}

func handleAPI(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	calName := r.URL.Query().Get("calName")

	switch r.Method {
	case "GET":
		password := getAuth(r)

		if password == "" {
			available, err := isURLAvailable(db, calName)
			if err != nil {
				fmt.Println("isURLAvailable error:", err)
				http.Error(w, "Server error", 500)
				return
			}
			if available {
				// unclaimed — frontend shows "claim" form
				w.WriteHeader(http.StatusNotFound)
				return
			}
			// exists but no password provided
			http.Error(w, "Unauthorized", 401)
			return
		}

		if !verifyPassword(db, calName, password) {
			http.Error(w, "Unauthorized", 401)
			return
		}

		startQueryWindow := r.URL.Query().Get("start")
		endQueryWindow := r.URL.Query().Get("end")
		entries, err := getCalendarEntries(db, calName, startQueryWindow, endQueryWindow)
		if err != nil {
			http.Error(w, "Error fetching calendar entries", 500)
			return
		}

		data, err := json.Marshal(entries)
		if err != nil {
			http.Error(w, "Error encoding calendar entries", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	case "POST":
		var body struct {
			ClaimingNewURL bool   `json:"claimingNewURL"`
			Password       string `json:"password"`
			CalEntry
		}
		json.NewDecoder(r.Body).Decode(&body)

		if body.ClaimingNewURL {
			handleClaimCalendar(db, w, calName, body.Password)
		} else {
			handleCreateEntry(db, w, r, calName, body.CalEntry)
		}

	case "PUT":
		password := getAuth(r)
		if !verifyPassword(db, calName, password) {
			http.Error(w, "Unauthorized", 401)
			return
		}

		var body struct {
			ID int `json:"id"`
			CalEntry
		}
		json.NewDecoder(r.Body).Decode(&body)

		err := updateCalendarEntry(db, body.ID, body.CalEntry)
		if err != nil {
			http.Error(w, "Error updating entry", 500)
			return
		}
		w.WriteHeader(200)

	case "DELETE":
		password := getAuth(r)
		if !verifyPassword(db, calName, password) {
			http.Error(w, "Unauthorized", 401)
			return
		}

		var body struct {
			ID int `json:"id"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		err := deleteCalendarEntry(db, body.ID)
		if err != nil {
			http.Error(w, "Error deleting entry", 500)
			return
		}
		w.WriteHeader(200)
	}
}

func handleClaimCalendar(db *sql.DB, w http.ResponseWriter, calName string, password string) {
	available, err := isURLAvailable(db, calName)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	if !available {
		http.Error(w, "Calendar name already taken", 409)
		return
	}
	err = savePassword(db, calName, password)
	if err != nil {
		http.Error(w, "Error creating calendar", 500)
		return
	}
	w.WriteHeader(201)
}

func handleCreateEntry(db *sql.DB, w http.ResponseWriter, r *http.Request, calName string, entry CalEntry) {
	password := getAuth(r)
	if !verifyPassword(db, calName, password) {
		http.Error(w, "Unauthorized", 401)
		return
	}
	err := createCalendarEntry(db, calName, entry)
	if err != nil {
		http.Error(w, "Error creating entry", 500)
		return
	}
	w.WriteHeader(201)
}
