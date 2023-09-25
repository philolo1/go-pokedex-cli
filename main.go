package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// hint that is not used yet
type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func help() error {
	fmt.Println(`
	Welcome to the Pokedex!
	Usage:

	help: Displays a help message
	exit: Exit the Pokedex`)

	return nil
}

func exit() error {
	os.Exit(0)
	return nil
}

func createMap(mapInfo *MapInfo) map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    help,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    exit,
		},
		"map": {
			name:        "next map",
			description: "Get next map",
			callback:    mapInfo.query,
		},
		"mapb": {
			name:        "previous map",
			description: "Get previous map",
			callback:    mapInfo.queryBack,
		},
	}
}

func main() {
	url := "https://pokeapi.co/api/v2/location-area"
	mapInfo := &MapInfo{
		currentUrl:  &url,
		previousUrl: nil,
	}

	cmdMap := createMap(mapInfo)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Pokedex >  ")
		text, _ := reader.ReadString('\n')

		text = strings.TrimSpace(text)

		el, ok := cmdMap[text]

		if ok {
			el.callback()
		} else {
			fmt.Println("You entered:", text)
		}

	}
}
