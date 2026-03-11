package main

import (
	"database/sql"
	"os"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "data/urlcal.db")
	if err != nil {
		os.Exit(1)
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS calendars (
		calendar_name TEXT PRIMARY KEY,
		password_hash TEXT NOT NULL
	)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		calendar_name TEXT NOT NULL,
		title TEXT NOT NULL,
		start TEXT NOT NULL,
		end TEXT NOT NULL,
		note TEXT,
		color TEXT,
		FOREIGN KEY(calendar_name) REFERENCES calendars(calendar_name)
	)`)

	return db
}

// URLs and Passwords

func isURLAvailable(db *sql.DB, calName string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM calendars WHERE calendar_name = ?)", calName).Scan(&exists)
	return !exists, err
}

func savePassword(db *sql.DB, calName string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO calendars (calendar_name, password_hash) VALUES (?, ?)", calName, string(hash))
	return err
}

func getPasswordHash(db *sql.DB, calName string) (string, error) {
	var hash string
	err := db.QueryRow("SELECT password_hash FROM calendars WHERE calendar_name = ?", calName).Scan(&hash)
	return hash, err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func verifyPassword(db *sql.DB, calName string, password string) bool {
	hash, err := getPasswordHash(db, calName)
	if err != nil {
		return false
	}
	return checkPasswordHash(password, hash)
}

// Calendar entries

func getCalendarEntries(db *sql.DB, calName string, startWindow string, endWindow string) ([]CalEntry, error) {
	rows, err := db.Query("SELECT id, title, start, end, note, color FROM entries WHERE calendar_name = ? AND start >= ? AND end <= ?", calName, startWindow, endWindow)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []CalEntry
	for rows.Next() {
		var entry CalEntry
		err := rows.Scan(&entry.ID, &entry.Title, &entry.Start, &entry.End, &entry.Note, &entry.Color)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func createCalendarEntry(db *sql.DB, calName string, entry CalEntry) error {
	_, err := db.Exec("INSERT INTO entries (calendar_name, title, start, end, note, color) VALUES (?, ?, ?, ?, ?, ?)",
		calName, entry.Title, entry.Start, entry.End, entry.Note, entry.Color)
	return err
}

func updateCalendarEntry(db *sql.DB, entryID int, entry CalEntry) error {
	_, err := db.Exec("UPDATE entries SET title = ?, start = ?, end = ?, note = ?, color = ? WHERE id = ?",
		entry.Title, entry.Start, entry.End, entry.Note, entry.Color, entryID)
	return err
}

func deleteCalendarEntry(db *sql.DB, entryID int) error {
	_, err := db.Exec("DELETE FROM entries WHERE id = ?", entryID)
	return err
}
