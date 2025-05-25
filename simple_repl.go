package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bilalhachim/pokedexcli/internal"
)

var caughtPokemon map[string]Pokemon = map[string]Pokemon{}
var cache internal.Cache = *internal.NewCache(300 * time.Millisecond)

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}
type location_areas struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
type PokemonEncountersInLocationArea struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}
type config struct {
	previousUrl string
	nextUrl     string
}

type cliCommand struct {
	config      *config
	name        string
	description string
	callback    func(*config, string) error
	parameter   string
}

func cleanInput(text string) []string {
	string_slice := strings.Split(text, " ")
	for i := 0; i < len(string_slice); i++ {
		if string_slice[i] == "" {
			string_slice = slices.Delete(string_slice, i, i+1)
			i = i - 1
		}
	}
	for i := 0; i < len(string_slice); i++ {
		string_slice[i] = strings.ToLower(string_slice[i])
	}
	for i := 0; i < len(string_slice); i++ {
		string_slice[i] = strings.Trim(string_slice[i], " ")
	}
	return string_slice
}
func clean_slice(slice_string []string) []string {
	var simple_slice []string
	for i := 0; i < len(slice_string); i++ {
		if slice_string[i] != "" {
			simple_slice = append(simple_slice, slice_string[i])
		}
	}
	return simple_slice
}
func get_commands() [][]string {
	var slices_string [][]string
	scan := bufio.NewScanner(os.Stdin)
	counter := 0
	for scan.Scan() {
		slice_string := make([]string, 1)
		slice_string = append(slice_string, cleanInput(scan.Text())...)
		slices_string = append(slices_string, clean_slice(slice_string))
		counter++

	}
	return slices_string
}
func simple_repl() {
	var config config
	var list_string [][]string
	commands := commandRegistry()
	for {

		fmt.Print("Pokedex >")
		list_string = get_commands()
		for i := 0; i < len(list_string); i++ {

			value, value_ok := commands[list_string[i][0]]
			if len(list_string[i]) > 1 {
				value.parameter = list_string[i][1]
			} else {
				value.parameter = ""
			}
			if !value_ok {
				fmt.Println("\nUnkown Command")
			} else {
				value.callback(&config, value.parameter)
			}

		}

	}
}
func commandExit(config *config, value_parameter string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(config *config, value_parameter string) error {
	fmt.Println("\nWelcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}
func commandMap(config *config, value_parameter string) error {

	locations, err := makeRequest(config)
	if err != nil {
		return err
	}
	for i := 0; i < len(locations.Results); i++ {
		fmt.Println(locations.Results[i].Name)
	}

	return nil
}
func commandRegistry() map[string]cliCommand {
	registry_commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displays the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "displays the pokemons in specific location areas the syntax is explore <area_name>",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "catch a specific pokemon the syntax is catch <pokemon_name>",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "prints the name, height, weight, stats and type(s) of the Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "prints a list of all the names of the Pokemon has caught",
			callback:    commandPokedex,
		},
	}
	return registry_commands
}
func commandPokedex(config *config, value_parameter string) error {
	if len(caughtPokemon) == 0 {
		fmt.Println("\nno Pokemon caught")
		return nil
	}
	fmt.Println("\nYour Pokedex:")
	for _, v := range caughtPokemon {
		fmt.Println(" - " + v.Name)
	}
	return nil
}
func commandInspect(config *config, value_parameter string) error {
	pokemon, ok := caughtPokemon[value_parameter]
	if !ok {
		fmt.Println("\nyou have not caught that pokemon")
		return nil
	}
	fmt.Println("Name: " + pokemon.Name)
	fmt.Println("Height: " + strconv.Itoa(pokemon.Height))
	fmt.Println("Weight: " + strconv.Itoa(pokemon.Weight))
	fmt.Println("Stats: ")
	for i := 0; i < len(pokemon.Stats); i++ {
		fmt.Println("  -" + pokemon.Stats[i].Stat.Name + ": " + strconv.Itoa(pokemon.Stats[i].BaseStat))
	}
	fmt.Println("Types: ")
	for i := 0; i < len(pokemon.Types); i++ {
		fmt.Println("  - " + pokemon.Types[i].Type.Name)
	}

	return nil
}
func commandCatch(config *config, value_parameter string) error {
	fmt.Println("\nThrowing a Pokeball at " + value_parameter + "...")
	pokemon, err := makeCatchRequest(value_parameter)
	if err != nil {
		return err
	}
	r := *rand.New(rand.NewSource(rand.Int63()))
	catchRate := r.Intn(100)
	if catchRate >= 50 {
		fmt.Println(pokemon.Name + " was caught!")
		fmt.Println("You may now inspect it with the inspect command.")
		caughtPokemon[pokemon.Name] = pokemon
	} else if catchRate < 50 {
		fmt.Println(pokemon.Name + " escaped!")

	}
	return nil
}
func commandExplore(config *config, value_parameter string) error {

	pokemon_encounters_in_location_area, err := makeExploreRequest(value_parameter)
	if err != nil {
		return err
	}
	fmt.Println("Exploring " + value_parameter)
	fmt.Println("Found Pokemon:")
	for i := 0; i < len(pokemon_encounters_in_location_area.PokemonEncounters); i++ {
		fmt.Println(" - " + pokemon_encounters_in_location_area.PokemonEncounters[i].Pokemon.Name)
	}
	return nil
}
func commandMapb(config *config, value_parameter string) error {
	locations, err := makeRequestBack(config)
	if err != nil {
		return err
	}
	for i := 0; i < len(locations.Results); i++ {
		fmt.Println(locations.Results[i].Name)
	}

	return nil
}
func makeRequestBack(config *config) (location_areas, error) {
	var temporary_locations location_areas
	if config.nextUrl == "" {
		rawUrl := "https://pokeapi.co/api/v2/location-area/"
		req, err := http.NewRequest(http.MethodGet, rawUrl, nil)
		if err != nil {
			return location_areas{}, err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return location_areas{}, err
		}
		resBody, err := io.ReadAll(res.Body)
		cache.Add(rawUrl, resBody)
		if err != nil {
			return location_areas{}, err
		}
		json.Unmarshal(resBody, &temporary_locations)
		config.nextUrl = temporary_locations.Next
		config.previousUrl = temporary_locations.Previous
		return temporary_locations, nil
	}
	value, exist := cache.Get(config.previousUrl)
	if !exist {
		req, err := http.NewRequest(http.MethodGet, config.previousUrl, nil)
		if err != nil {
			return location_areas{}, err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return location_areas{}, err
		}
		resBody, err := io.ReadAll(res.Body)
		cache.Add(config.previousUrl, resBody)
		if err != nil {
			return location_areas{}, err
		}
		json.Unmarshal(resBody, &temporary_locations)
		config.nextUrl = temporary_locations.Next
		config.previousUrl = temporary_locations.Previous
		return temporary_locations, nil
	}
	json.Unmarshal(value, &temporary_locations)
	config.nextUrl = temporary_locations.Next
	config.previousUrl = temporary_locations.Previous
	return temporary_locations, nil
}
func makeRequest(config *config) (location_areas, error) {
	var temporary_locations location_areas
	if config.nextUrl == "" {
		rawUrl := "https://pokeapi.co/api/v2/location-area/"
		req, err := http.NewRequest(http.MethodGet, rawUrl, nil)
		if err != nil {
			return location_areas{}, err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return location_areas{}, err
		}
		resBody, err := io.ReadAll(res.Body)
		cache.Add(rawUrl, resBody)
		if err != nil {
			return location_areas{}, err
		}
		json.Unmarshal(resBody, &temporary_locations)
		config.nextUrl = temporary_locations.Next
		config.previousUrl = temporary_locations.Previous
		return temporary_locations, nil
	}
	value, exist := cache.Get(config.nextUrl)
	if !exist {
		req, err := http.NewRequest(http.MethodGet, config.nextUrl, nil)
		if err != nil {
			return location_areas{}, err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return location_areas{}, err
		}
		resBody, err := io.ReadAll(res.Body)
		cache.Add(config.nextUrl, resBody)
		if err != nil {
			return location_areas{}, err
		}
		json.Unmarshal(resBody, &temporary_locations)
		config.nextUrl = temporary_locations.Next
		config.previousUrl = temporary_locations.Previous
		return temporary_locations, nil
	}

	json.Unmarshal(value, &temporary_locations)
	config.nextUrl = temporary_locations.Next
	config.previousUrl = temporary_locations.Previous
	return temporary_locations, nil

}
func makeExploreRequest(location_area string) (PokemonEncountersInLocationArea, error) {
	var Pokemons_encounter PokemonEncountersInLocationArea
	rawUrl := "https://pokeapi.co/api/v2/location-area/" + location_area
	if location_area == "" {
		return PokemonEncountersInLocationArea{}, nil
	}
	jsonPokemon, is_exist := cache.Get(rawUrl)
	if is_exist {
		json.Unmarshal(jsonPokemon, &Pokemons_encounter)
		return Pokemons_encounter, nil
	}
	req, err := http.NewRequest(http.MethodGet, rawUrl, nil)
	if err != nil {
		return PokemonEncountersInLocationArea{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return PokemonEncountersInLocationArea{}, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return PokemonEncountersInLocationArea{}, err
	}
	cache.Add(rawUrl, resBody)
	json.Unmarshal(resBody, &Pokemons_encounter)
	return Pokemons_encounter, nil
}

func makeCatchRequest(pokemon_name string) (Pokemon, error) {
	var pokemon Pokemon
	rawUrl := "https://pokeapi.co/api/v2/pokemon/" + pokemon_name
	if pokemon_name == "" {
		return Pokemon{}, nil
	}
	pokemon, is_exist := caughtPokemon[pokemon_name]
	if is_exist {
		return pokemon, nil
	}
	req, err := http.NewRequest(http.MethodGet, rawUrl, nil)
	if err != nil {
		return Pokemon{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}
	cache.Add(rawUrl, resBody)
	json.Unmarshal(resBody, &pokemon)
	return pokemon, nil
}
