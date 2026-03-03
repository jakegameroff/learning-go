package main

import (
	"fmt"
	"math/rand"
)

func startGame() {
	secretNumber := rand.Intn(100) + 1
	var guess int
    for guess != secretNumber {
		fmt.Print("What is your guess: ")
		fmt.Scan(&guess)

		if guess < secretNumber {
			fmt.Print("Too low!\n")
		} else if guess > secretNumber {
			fmt.Print("Too high!\n")
		} else {
			fmt.Printf("Congratulations! Your guess of %d is correct! :)\n", secretNumber)
		}
	}
}

func main() {
	startGame()
}
