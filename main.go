package main

import (
	"bufio"
	"fmt"
	"github.com/Nic-Vladimir/pokedexcli/internal"
	"os"
	"strings"
	"time"
)

type config struct {
	next     string
	previous string
	pokedex  map[string]PokemonData
	cache    *internal.Cache
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	tokens := strings.Fields(text)
	return tokens
}

// --- Main ---
func main() {
	config := config{
		pokedex: make(map[string]PokemonData),
		cache:   internal.NewCache(5 * time.Minute),
	}
	input := bufio.NewScanner(os.Stdin)
	commands := GetCommands()
	for {
		fmt.Print("Pokedex > ")
		input.Scan()

		tokens := cleanInput(input.Text())
		if len(tokens) == 0 {
			continue
		}

		userCommand := tokens[0]
		param := ""
		if len(tokens) > 1 {
			param = tokens[1]
		}

		if cmd, ok := commands[userCommand]; ok {
			if err := cmd.callback(&config, param); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
