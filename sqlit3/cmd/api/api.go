package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (app *app) mount() http.Handler {
	r := chi.NewRouter() // TODO: add middleware
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health check passed!"))
	})

	return r
}

func (app *app) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Minute * 5,
		ReadTimeout:  time.Minute,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server on %s", app.config.addr)

	return srv.ListenAndServe()
}

type app struct {
	config config
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
