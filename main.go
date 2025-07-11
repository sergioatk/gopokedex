package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sergioatk/gopokedex/internal/pokecache"
	"github.com/sergioatk/gopokedex/internal/pokemon"
)

var commandList map[string]cliCommand

func main() {

	commandList = map[string]cliCommand{
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
			description: "Displays pokemon locations. Each time you call this endpoint, a new page is requested.",
			callback:    pokemon.Map,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous pokemon locations. Each time you call this endpoint, a new page is requested.",
			callback:    pokemon.Mapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore the different pokemons that inhabit an area.",
			callback:    pokemon.Explore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon.",
			callback:    pokemon.Catch,
		},
		"inspect": {
			name:        "inspect",
			description: "View stats about your caught pokemon.",
			callback:    pokemon.Inspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all your caught pokemon.",
			callback:    pokemon.Pokedex,
		},
	}

	config := &pokemon.CommandConfig{}

	cache := pokecache.NewCache(60 * time.Second)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			rawUserInput := scanner.Text()
			cleanedInput := cleanInput(rawUserInput)
			desiredCommand := cleanedInput[0]
			parameter := ""
			if len(cleanedInput) >= 2 {
				parameter = cleanedInput[1]
			}

			command, ok := commandList[desiredCommand]
			if !ok {
				fmt.Println("Unknown command")
				commandExit(config, cache, "")
			}

			err := command.callback(config, cache, parameter)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}

}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	return strings.Fields(lowered)
}

func commandExit(config *pokemon.CommandConfig, cache *pokecache.Cache, parameter string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *pokemon.CommandConfig, cache *pokecache.Cache, parameter string) error
}

func commandHelp(config *pokemon.CommandConfig, cache *pokecache.Cache, parameter string) error {

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, metadata := range commandList {
		fmt.Printf("%s: %s\n", metadata.name, metadata.description)
	}

	return nil

}
