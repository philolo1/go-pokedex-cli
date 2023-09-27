package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/philolo1/go-pokedex-cli/cache"
	. "github.com/philolo1/go-pokedex-cli/cache"
)

type MapInfo struct {
	previousUrl *string
	currentUrl  *string
	cache       Cache
}

type ExploreInfo struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func NewMapInfo() *MapInfo {
	url := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	duration := time.Second * 60
	return &MapInfo{
		currentUrl:  &url,
		previousUrl: nil,
		cache:       cache.NewCache(duration),
	}
}

type MapResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous *string  `json:"previous"` // *string since it can be null
	Results  []Result `json:"results"`
}

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func getResponse[Response any](url string, cache *Cache) (Response, error) {

	// fmt.Printf("checking cache for %v\n", url)

	item, hasCache := cache.Get(url)

	var result Response

	if hasCache {
		println("Found cache")

		err := json.Unmarshal(item, &result)

		if err == nil {
			println("Return cached result")
			return result, nil
		} else {
			fmt.Printf("Cache not found %v", err)
		}
	}

	// fmt.Printf("doing the query %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}

	if resp.StatusCode >= 300 {
		return result, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	defer resp.Body.Close()

	// Unmarshal the response into a map
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	go func() {
		bytes, err := json.Marshal(result)

		if err == nil {
			cache.Add(url, bytes)
		}
	}()

	return result, nil
}

func (m *MapInfo) Query(params *[]string) error {

	if m.currentUrl == nil || len(*m.currentUrl) == 0 {
		return errors.New("No next url")
	}
	result, err := getResponse[MapResponse](*m.currentUrl, &m.cache)

	if err != nil {
		return err
	}

	for _, item := range result.Results {
		fmt.Printf("%v\n", item.Name)
	}

	m.previousUrl = m.currentUrl
	m.currentUrl = &result.Next

	return nil
}

func (m *MapInfo) QueryBack(params *[]string) error {
	if m.previousUrl == nil || len(*m.previousUrl) == 0 {
		return errors.New("No previous url")
	}
	result, err := getResponse[MapResponse](*m.previousUrl, &m.cache)

	if err != nil {
		return err
	}

	for _, item := range result.Results {
		fmt.Printf("%v\n", item.Name)
	}

	m.currentUrl = m.previousUrl
	m.previousUrl = result.Previous

	return nil
}

func (m *MapInfo) ExploreRegion(params *[]string) error {
	if len(*params) != 1 {
		return errors.New("Exactly one parameter required")
	}

	location := (*params)[0]

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location)

	fmt.Printf("Exploring %v...\n", location)

	result, err := getResponse[ExploreInfo](url, &m.cache)

	if err != nil {
		return err
	}

	fmt.Printf("Found Pokemon:\n")
	for _, pokemon := range result.PokemonEncounters {
		fmt.Println("- " + pokemon.Pokemon.Name)
	}

	return nil
}
