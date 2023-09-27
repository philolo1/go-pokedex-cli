package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/philolo1/go-pokedex-cli/api"
)

// hint that is not used yet
type cliCommand struct {
	name        string
	description string
	callback    func(params *[]string) error
}

func help(params *[]string) error {
	fmt.Println(`
	Welcome to the Pokedex!
	Usage:
	help: Displays a help message
	exit: Exit the Pokedex
	map: Get next map
	mapb: Get previous map
	explore <region>: List pokemon in region
	catch <pokemon>: Catch pokemon
	inspect <pokemon>: Inspect pokemon in pokedex
	pokedex: Show current pokedex
	`)

	return nil
}

func exit(params *[]string) error {
	os.Exit(0)
	return nil
}

func createMap(mapInfo *api.MapInfo) map[string]cliCommand {
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
			callback:    mapInfo.Query,
		},
		"mapb": {
			name:        "previous map",
			description: "Get previous map",
			callback:    mapInfo.QueryBack,
		},
		"explore": {
			name:        "explore region",
			description: "list pokemon in region",
			callback:    mapInfo.ExploreRegion,
		},
		"catch": {
			name:        "catch pokemon",
			description: "catch pokemon",
			callback:    mapInfo.CatchPokemon,
		},
		"inspect": {
			name:        "inspect pokemon in pokedex",
			description: "inspect pokemon in pokedex",
			callback:    mapInfo.InspectPokemon,
		},
		"pokedex": {
			name:        "show current pokedex",
			description: "show current pokedex",
			callback:    mapInfo.ShowPokedex,
		},
	}
}

func main() {
	mapInfo := api.NewMapInfo()

	cmdMap := createMap(mapInfo)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Pokedex >  ")
		text, _ := reader.ReadString('\n')

		wordArr := strings.Fields(text)

		text = wordArr[0]

		param := wordArr[1:]

		el, ok := cmdMap[text]

		if ok {
			err := el.callback(&param)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			}
		} else {
			fmt.Println("You entered:", text)
		}

	}
}
