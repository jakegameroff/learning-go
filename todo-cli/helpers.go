package main

import (
	"fmt"
	"os"
	"strconv"
	"encoding/json"
)

func save(todoList []Todo) {
	data, err := json.Marshal(todoList)
	if err != nil {
		fmt.Print("Error converting to json")
		os.Exit(1)
	}
	err = os.WriteFile("todolist.json", data, 0644)
	if err != nil {
		fmt.Print("Error writing to file")
		os.Exit(1)
	}
}

func fetch() []Todo {
	var todoList []Todo
	data, err := os.ReadFile("todolist.json")

	if err != nil {
		err = os.WriteFile("todolist.json", []byte(""), 0644)
		if err != nil {
			fmt.Print("Error fetching file")
			os.Exit(1)
		} 
		return todoList
	}

	err = json.Unmarshal(data, &todoList)
	if err != nil {
		fmt.Print("Error unpacking json")
		os.Exit(1)
	}
	return todoList
}

func stringifyTodo(todo Todo) string {
	var completionStatus string
	var todoString string

	if todo.Done {
		completionStatus = "[x] "
	} else {
		completionStatus = "[] "
	}
	todoString = completionStatus + todo.Text + "\n"
	return todoString
}

func readTodoList(todoList []Todo) {
	fmt.Print("Printing todo list\n******************\n")
	for i, item := range todoList {
		indexString := strconv.Itoa(i) + ". "
		todoString := indexString + stringifyTodo(item)
		fmt.Print(todoString)
	}
}

func addTodo(todoList *[]Todo, text string) {
	todo := Todo{Text: text, Done: false}
	*todoList = append(*todoList, todo)
}

func deleteTodo(todoList *[]Todo, index int) Todo {
	if index < 0 || index >= len(*todoList) {
		fmt.Print("Index out of range")
		os.Exit(1)
	}
	todo := (*todoList)[index]

	sliceBefore := (*todoList)[:index]
	sliceAfter := (*todoList)[index+1:]
	newTodoList := append(sliceBefore, sliceAfter...)
	*todoList = newTodoList

	return todo
}

func modifyTodo(todoList *[]Todo, mark bool, index int) *Todo {
	if index < 0 || index >= len(*todoList) {
		fmt.Print("Index out of range")
		os.Exit(1)
	}
	todo := &(*todoList)[index]
	todo.Done = mark

	return todo
}