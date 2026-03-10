package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func save(todoList []Todo) error {
	data, err := json.Marshal(todoList)
	if err != nil {
		return err
	}
	err = os.WriteFile("todolist.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func fetch() ([]Todo, error) {
	var todoList []Todo
	data, err := os.ReadFile("todolist.json")

	if err != nil {
		err = os.WriteFile("todolist.json", []byte("[]"), 0644)
		if err != nil {
			return todoList, err
		}
		return todoList, nil
	}

	err = json.Unmarshal(data, &todoList)
	if err != nil {
		return todoList, err
	}
	return todoList, nil
}

func addTodo(todoList *[]Todo, text string) {
	todo := Todo{Text: text, Done: false}
	*todoList = append(*todoList, todo)
}

func deleteTodo(todoList *[]Todo, index int) (Todo, error) {
	if index < 0 || index >= len(*todoList) {
		return Todo{}, fmt.Errorf("Index out of range")
	}
	todo := (*todoList)[index]

	sliceBefore := (*todoList)[:index]
	sliceAfter := (*todoList)[index+1:]
	newTodoList := append(sliceBefore, sliceAfter...)
	*todoList = newTodoList

	return todo, nil
}

func modifyTodo(todoList *[]Todo, mark bool, index int) (*Todo, error) {
	if index < 0 || index >= len(*todoList) {
		return &Todo{}, fmt.Errorf("Index out of range")
	}
	todo := &(*todoList)[index]
	(*todo).Done = mark
	return todo, nil
}
