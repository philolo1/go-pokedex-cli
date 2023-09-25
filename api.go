package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

func (m *MapInfo) query() error {

	println("doing the query")

	result := getResponse(*m.currentUrl)

	for _, item := range result.Results {
		fmt.Printf("%v\n", item.Name)
	}

	m.previousUrl = m.currentUrl
	m.currentUrl = &result.Next

	return nil
}

func (m *MapInfo) queryBack() error {
	result := getResponse(*m.previousUrl)

	for _, item := range result.Results {
		fmt.Printf("%v\n", item.Name)
	}

	m.currentUrl = m.previousUrl
	m.previousUrl = result.Previous

	return nil
}
