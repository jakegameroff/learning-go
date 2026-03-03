package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

func startGame() {
	numGuesses := 0
	secretNumber := rand.Intn(100) + 1
	
	var guess int
    for guess != secretNumber {
		fmt.Print("What is your guess: ")
		
		var input string
		var err error

		fmt.Scan(&input)
		guess, err = strconv.Atoi(input)

		if err != nil {
			fmt.Print("Please enter a number!\n")
			continue
		}

		numGuesses += 1

		if guess < secretNumber {
			fmt.Print("Too low!\n")
		} else if guess > secretNumber {
			fmt.Print("Too high!\n")
		} else {
			fmt.Printf(
				"Congratulations! Your guess of %d is correct! It took %d guesses.",
				secretNumber, numGuesses)
		}
	}
}

func main() {
	startGame()
}
