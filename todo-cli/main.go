package main

import (
	"fmt"
	"bufio"
	"os"
)

type Todo struct {
	Text string
	Done bool
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	todoList := fetch()

	for {
		var action string
		fmt.Print("What would you like to do? (read r, add a, delete d, mark m, unmark u, exit e): ")
		fmt.Scan(&action)

		switch action {
		case "r":
			readTodoList(todoList)
		case "a":
			fmt.Print("Task to add: ")
			scanner.Scan()
			text := scanner.Text()
			addTodo(&todoList, text)
			save(todoList)
		case "d":
			var index int
			readTodoList(todoList)
			fmt.Print("Index to delete: ")
			fmt.Scan(&index)
			deleteTodo(&todoList, index)
			save(todoList)
		case "m":
			var index int
			readTodoList(todoList)
			fmt.Print("Index to mark: ")
			fmt.Scan(&index)
			modifyTodo(&todoList, true, index)
			save(todoList)
		case "u":
			var index int
			readTodoList(todoList)
			fmt.Print("Index to unmark: ")
			fmt.Scan(&index)
			modifyTodo(&todoList, false, index)
			save(todoList)
		case "e":
			fmt.Print("Bye!")
			os.Exit(0)
		default:
			fmt.Print("Unknown command")
		}
	}
}