package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type MapInfo struct {
	previousUrl *string
	currentUrl  *string
	cache       Cache
}

func NewMapInfo() *MapInfo {
	url := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	duration := time.Second * 60
	return &MapInfo{
		currentUrl:  &url,
		previousUrl: nil,
		cache:       NewCache(duration),
	}
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

func getResponse(url string, cache *Cache) Response {

	fmt.Printf("checking cache for %v\n", url)

	item, hasCache := cache.Get(url)

	var result Response

	if hasCache {
		println("Found cache")

		err := json.Unmarshal(item, &result)

		if err == nil {
			println("Return cached result")
			return result
		} else {
			fmt.Printf("Cache not found %v", err)
		}
	}

	fmt.Printf("doing the query %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Unmarshal the response into a map
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	go func() {
		bytes, err := json.Marshal(result)

		if err == nil {
			cache.Add(url, bytes)
		}
	}()

	return result
}

func (m *MapInfo) query() error {

	if m.currentUrl == nil || len(*m.currentUrl) == 0 {
		return errors.New("No next url")
	}
	result := getResponse(*m.currentUrl, &m.cache)

	for _, item := range result.Results {
		fmt.Printf("%v\n", item.Name)
	}

	m.previousUrl = m.currentUrl
	m.currentUrl = &result.Next

	return nil
}

func (m *MapInfo) queryBack() error {
	if m.previousUrl == nil || len(*m.previousUrl) == 0 {
		return errors.New("No previous url")
	}
	result := getResponse(*m.previousUrl, &m.cache)

	for _, item := range result.Results {
		fmt.Printf("%v\n", item.Name)
	}

	m.currentUrl = m.previousUrl
	m.previousUrl = result.Previous

	return nil
}
