package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/philolo1/go-pokedex-cli/cache"
	. "github.com/philolo1/go-pokedex-cli/cache"
)

type MapInfo struct {
	previousUrl *string
	currentUrl  *string
	cache       Cache
	pokedex     map[string]PokemonInfo
}

type PokemonInfo struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience         int    `json:"base_experience"`
	Height                 int    `json:"height"`
	HeldItems              []any  `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	PastTypes []any  `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func (p PokemonInfo) String() string {

	var statBuilder strings.Builder

	for _, item := range p.Stats {
		statBuilder.WriteString(fmt.Sprintf("  - %v: %v\n", item.Stat.Name, item.BaseStat))
	}

	var typesBuilder strings.Builder

	for _, item := range p.Types {
		typesBuilder.WriteString(fmt.Sprintf("  - %v\n", item.Type.Name))
	}

	return fmt.Sprintf(
		`
Name: %s
Height: %v
Weight: %v
Stats: 
%vTypes: 
%v
`, p.Name, p.Height, p.Weight, statBuilder.String(), typesBuilder.String())

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
		pokedex:     make(map[string]PokemonInfo),
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

func hasCaught(prob int, exp int) bool {
	if exp > 100 {
		return prob > 5
	} else {
		return prob > 2
	}
}

func (m *MapInfo) CatchPokemon(params *[]string) error {
	if len(*params) != 1 {
		return errors.New("Exactly one parameter required")
	}

	pokemon := (*params)[0]

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon)

	result, err := getResponse[PokemonInfo](url, &m.cache)

	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon)

	exp := result.BaseExperience

	prob := rand.Intn(10) + 1

	if hasCaught(prob, exp) {
		fmt.Printf("%v was caught!\n", pokemon)
		m.pokedex[pokemon] = result
	} else {
		fmt.Printf("%v escaped!\n", pokemon)
	}

	return nil
}

func (m *MapInfo) InspectPokemon(params *[]string) error {
	if len(*params) != 1 {
		return errors.New("Exactly one parameter required")
	}

	pokemonName := (*params)[0]

	pokemon, ok := m.pokedex[pokemonName]

	if !ok {
		return errors.New("you have not caught that pokemon")
	}

	fmt.Printf("%v", pokemon)

	return nil
}
