package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// hint that is not used yet
// type cliCommand struct {
// 	name        string
// 	description string
// 	callback    func() error
// }
//
// return map[string]cliCommand{
//     "help": {
//         name:        "help",
//         description: "Displays a help message",
//         callback:    commandHelp,
//     },
//     "exit": {
//         name:        "exit",
//         description: "Exit the Pokedex",
//         callback:    commandExit,
//     },
// }

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
