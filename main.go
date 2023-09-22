package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("pokedex >  ")
		text, _ := reader.ReadString('\n')

		text = strings.TrimSpace(text)

		switch text {
		case "exit":
			return
		case "help":
			fmt.Println(`
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex`)

		default:

			fmt.Println("You entered:", text)
		}

	}
}
