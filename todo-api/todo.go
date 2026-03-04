package main

import (
	"net/http"
	"encoding/json"
	"strconv"
	"strings"
)

type Todo struct {
	Text string
	Done bool
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	todoList, err := fetch()

	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}

	switch r.Method {
	case "GET":
		data, err := json.Marshal(todoList)
		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	case "POST":
		var body map[string]string
		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.Decode(&body)

		text := body["text"]
		addTodo(&todoList, text)
		err := save(todoList)

		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}

		w.WriteHeader(http.StatusCreated)

	case "PUT":
		var body map[string]string
		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.Decode(&body)

		index, err := strconv.Atoi(body["index"])
		mark := strings.ToLower(body["mark"]) == "true"

		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}

		_, err = modifyTodo(&todoList, mark, index)

		if err != nil {
			http.Error(w, "Bad request", 400)
			return
		}

		err = save(todoList)

		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}

		w.WriteHeader(http.StatusOK)

	case "DELETE":
		var body map[string]int
		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.Decode(&body)

		index := body["index"]
		_, err := deleteTodo(&todoList, index)

		if err != nil {
			http.Error(w, "Bad request", 400)
			return
		}

		err = save(todoList)

		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}