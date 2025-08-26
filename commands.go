package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
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
			description: "Displays the next 20 locations in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations in the Pokemon world",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "explore <area_name> - Returns a list of Pokemon located there",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "catch <pokemon_name> - Attempt to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect <pokemon_name> - Displays information about a previously caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays a list of caught Pokemon",
			callback:    commandPokedex,
		},
	}
}

// --- Commands ---

func commandPokedex(cfg *config, param string) error {
	fmt.Println("Your Pokedex:")
	for _, p := range cfg.pokedex {
		fmt.Printf(" - %s\n", p.Name)
	}
	return nil
}

func commandInspect(cfg *config, pokemonName string) error {
	pokemonData, ok := cfg.pokedex[pokemonName]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name:   %s\n", pokemonData.Name)
	fmt.Printf("Height: %d\n", pokemonData.Height)
	fmt.Printf("Weight: %d\n", pokemonData.Weight)

	fmt.Println("Stats:")
	for _, s := range pokemonData.Stats {
		fmt.Printf("  - %s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range pokemonData.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}

func commandCatch(cfg *config, pokemonName string) error {
	fmt.Println("Throwing a Pokeball at " + pokemonName + "...")
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName + "/"
	bodyBytes, err := GetRequest(url)
	var pokemonData PokemonData
	if err != nil {
		return fmt.Errorf("Catch: Invalid Pokemon name\nDetails: %s", err)
	}
	if err = json.Unmarshal(bodyBytes, &pokemonData); err != nil {
		return err
	}
	playerRoll := rand.Intn(6)
	pokemonRoll := pokemonData.BaseExperience / 50
	if playerRoll >= pokemonRoll {
		fmt.Println(pokemonName, "was caught!")
		cfg.pokedex[pokemonName] = pokemonData
	} else {
		fmt.Println(pokemonName, "escaped!")
	}
	return nil
}

func commandExplore(cfg *config, area string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + area + "/"
	fmt.Println("Exploring ", area, "...")
	var areaData AreaData
	if val, ok := cfg.cache.Get(url); ok {
		if err := json.Unmarshal(val, &areaData); err != nil {
			return err
		}
	} else {
		bodyBytes, err := GetRequest(url)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, bodyBytes)
		if err = json.Unmarshal(bodyBytes, &areaData); err != nil {
			return err
		}
	}
	fmt.Println("Found Pokemon: ")
	for _, p := range areaData.PokemonEncounters {
		fmt.Println(" - ", p.Pokemon.Name)
	}
	return nil
}

func commandMapB(cfg *config, param string) error {
	var url string
	if cfg.previous != "" {
		url = cfg.previous
	} else {
		return fmt.Errorf("No previous locations")
	}
	var resBody Locations
	if val, ok := cfg.cache.Get(url); ok {
		if err := json.Unmarshal(val, &resBody); err != nil {
			return err
		}
		fmt.Println("Using cached data")
	} else {
		bodyBytes, err := GetRequest(url)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, bodyBytes)
		if err = json.Unmarshal(bodyBytes, &resBody); err != nil {
			return err
		}
	}
	for _, r := range resBody.Results {
		fmt.Println(r.Name)
	}
	cfg.next = resBody.Next
	cfg.previous = resBody.Previous

	return nil
}

func commandMap(cfg *config, param string) error {
	var url string
	if cfg.next != "" {
		url = cfg.next
	} else {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	var resBody Locations
	if val, ok := cfg.cache.Get(url); ok {
		if err := json.Unmarshal(val, &resBody); err != nil {
			return err
		}
		fmt.Println("Using cached data")
	} else {
		bodyBytes, err := GetRequest(url)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, bodyBytes)
		if err = json.Unmarshal(bodyBytes, &resBody); err != nil {
			return err
		}
	}
	for _, r := range resBody.Results {
		fmt.Println(r.Name)
	}
	cfg.next = resBody.Next
	cfg.previous = resBody.Previous
	return nil
}

func commandExit(conf *config, param string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, param string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range GetCommands() {
		fmt.Printf("%-*s %s\n", 10, cmd.name+":", cmd.description)
	}
	return nil
}
