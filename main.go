package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

type MapInfo struct {
	previousUrl *string
	currentUrl  *string
}

type Response struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous *string  `json:"previous"` // *string since it can be null
	Results  []Result `json:"results"`
}

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func getResponse(url string) Response {
	fmt.Printf("doing the query %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Unmarshal the response into a map
	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	return result
}

func (m *MapInfo) query() {

	println("doing the query")

	result := getResponse(*m.currentUrl)

	for _, item := range result.Results {
		fmt.Printf("%v\n", item.Name)
	}

	m.previousUrl = m.currentUrl
	m.currentUrl = &result.Next
}

func (m *MapInfo) queryBack() {
	result := getResponse(*m.previousUrl)

	for _, item := range result.Results {
		fmt.Printf("%v\n", item.Name)
	}

	m.currentUrl = m.previousUrl
	m.previousUrl = result.Previous
}

func main() {
	url := "https://pokeapi.co/api/v2/location-area"
	mapInfo := &MapInfo{
		currentUrl:  &url,
		previousUrl: nil,
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Pokedex >  ")
		text, _ := reader.ReadString('\n')

		text = strings.TrimSpace(text)

		switch text {
		case "exit":
			return
		case "map":
			mapInfo.query()
			break
		case "mapb":
			mapInfo.queryBack()
			break
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
