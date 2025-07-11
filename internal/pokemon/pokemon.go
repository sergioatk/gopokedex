package pokemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"

	pokecache "github.com/sergioatk/gopokedex/internal/pokecache"
)

type ExploreResponse struct {
	AreaName          string `json:"area_name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type CommandConfig struct {
	Next     string
	Previous string
}

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ApiResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Result `json:"results"`
}

func getDataAndSetCache(url string, cache *pokecache.Cache) ([]byte, error) {
	var result []byte
	res, err := http.Get(url)
	if err != nil {
		return result, errors.New("error while exploring area")
	}

	defer res.Body.Close()
	if res.StatusCode == 404 {
		return result, errors.New("it appears that does not relly exist... no te quieras pasar de listo gatito :3")
	}
	if res.StatusCode != 200 {
		return result, errors.New("requested data has failed")

	}
	result, err = io.ReadAll(res.Body)
	if err != nil {
		return result, errors.New("error while exploring area")
	}
	cache.Add(url, result)

	return result, nil
}

func Explore(config *CommandConfig, cache *pokecache.Cache, parameter string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + parameter

	cachedValue, cacheHit := cache.Get(url)
	var apiResponse ExploreResponse
	var rawData []byte
	if cacheHit {
		rawData = cachedValue
	} else {

		apiRawData, err := getDataAndSetCache(url, cache)
		rawData = apiRawData

		if err != nil {
			return err
		}

	}

	if err := json.Unmarshal(rawData, &apiResponse); err != nil {
		return errors.New("error while unmarshaling response")
	}

	fmt.Printf("Exploring %s...\n", parameter)
	fmt.Println("Found Pokemon:")

	for _, pokemon := range apiResponse.PokemonEncounters {
		fmt.Printf("- %s\n", pokemon.Pokemon.Name)
	}

	return nil
}

func Map(config *CommandConfig, cache *pokecache.Cache, parameter string) error {
	url := "https://pokeapi.co/api/v2/location-area"

	if config.Next != "" {
		url = config.Next
	}

	cachedValue, cacheHit := cache.Get(url)
	var apiResponse ApiResponse
	var rawData []byte
	if cacheHit {
		rawData = cachedValue
	} else {

		apiRawData, err := getDataAndSetCache(url, cache)
		rawData = apiRawData

		if err != nil {
			return err
		}
	}

	if err := json.Unmarshal(rawData, &apiResponse); err != nil {
		return errors.New("error while unmarshaling response")
	}

	config.Next = apiResponse.Next
	config.Previous = apiResponse.Previous

	for _, result := range apiResponse.Results {
		fmt.Println(result.Name)
	}

	return nil

}

func Mapb(config *CommandConfig, cache *pokecache.Cache, parameter string) error {

	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	url := config.Previous

	cachedValue, cacheHit := cache.Get(url)
	var apiResponse ApiResponse
	var rawData []byte
	if cacheHit {
		rawData = cachedValue
	} else {
		apiRawData, err := getDataAndSetCache(url, cache)
		rawData = apiRawData

		if err != nil {
			return err
		}
	}

	if err := json.Unmarshal(rawData, &apiResponse); err != nil {
		return errors.New("error while unmarshaling response")
	}

	config.Next = apiResponse.Next
	config.Previous = apiResponse.Previous

	for _, result := range apiResponse.Results {
		fmt.Println(result.Name)
	}

	return nil

}

type PokemonDetail struct {
	BaseExperience int    `json:"base_experience"`
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			StatName string `json:"stat_name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			TypeName string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func isPokemonCaught(pokemon string, cache *pokecache.Cache) bool {
	_, pokemonIsCaught := cache.Get(pokemon)

	return pokemonIsCaught
}

func validateParameter(parameter string) error {
	if parameter == "" {
		return errors.New("please pass a valid param :/")
	}

	return nil
}

func Catch(config *CommandConfig, cache *pokecache.Cache, parameter string) error {

	paramErr := validateParameter(parameter)
	if paramErr != nil {
		return paramErr
	}

	pokemonAlreadyCaught := isPokemonCaught(parameter, cache)

	if pokemonAlreadyCaught {
		fmt.Println("You already caught this pokemon.")
		fmt.Println("Try to catch another one!")
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", parameter)

	url := "https://pokeapi.co/api/v2/pokemon/" + parameter

	cachedValue, cacheHit := cache.Get(url)
	var apiResponse PokemonDetail
	var rawData []byte
	if cacheHit {
		rawData = cachedValue
	} else {
		apiRawData, err := getDataAndSetCache(url, cache)
		rawData = apiRawData

		if err != nil {
			return err
		}
	}

	if err := json.Unmarshal(rawData, &apiResponse); err != nil {
		fmt.Println("error", err)
		return errors.New("error while unmarshaling response")
	}

	chance := math.Max(5, 100-(float64(apiResponse.BaseExperience)/200)*95)
	roll := rand.Float64() * 100
	caught := roll <= chance

	if caught {
		fmt.Printf("%s was caught!\n", parameter)
		// this means the pokemon was caught
		cache.Add(parameter, []byte{})
		err := addToPokedex(cache, parameter)
		if err != nil {

			return err
		}
	} else {
		fmt.Printf("%s escaped!\n", parameter)
	}

	return nil

}

func addToPokedex(cache *pokecache.Cache, parameter string) error {

	pokedex, pokedexInitialized := cache.Get("pokedex")
	if !pokedexInitialized {
		value, err := json.Marshal([]string{})

		if err != nil {

			return err
		}
		cache.Add("pokedex", value)

		pokedex = value
	}

	var storedPokedex []string

	err := json.Unmarshal(pokedex, &storedPokedex)

	if err != nil {
		return err
	}

	storedPokedex = append(storedPokedex, parameter)

	value, err := json.Marshal(storedPokedex)
	if err != nil {
		return err
	}
	cache.Add("pokedex", value)

	return nil

}

func Inspect(config *CommandConfig, cache *pokecache.Cache, parameter string) error {
	paramErr := validateParameter(parameter)
	if paramErr != nil {
		return paramErr
	}
	pokemonAlreadyCaught := isPokemonCaught(parameter, cache)

	if !pokemonAlreadyCaught {
		fmt.Println("It appears you haven't caught this pokemon yet, go for it!")
		return nil
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + parameter

	cachedValue, cacheHit := cache.Get(url)
	var apiResponse PokemonDetail
	var rawData []byte
	if cacheHit {
		rawData = cachedValue
	} else {
		fmt.Println("Sorry there was a problem with the cached data")
		return nil
	}

	if err := json.Unmarshal(rawData, &apiResponse); err != nil {
		return errors.New("error while unmarshaling response")
	}

	fmt.Printf("Name: %s\n", apiResponse.Name)
	fmt.Printf("Height: %s\n", apiResponse.Name)
	fmt.Printf("Weight: %s\n", apiResponse.Name)
	fmt.Printf("Stats: %s\n", apiResponse.Name)

	for _, stat := range apiResponse.Stats {
		fmt.Printf(" -%s: %d\n", stat.Stat.StatName, stat.BaseStat)
	}

	fmt.Printf("Types: %s\n", apiResponse.Name)

	for _, pokeType := range apiResponse.Types {
		fmt.Printf(" - %s\n", pokeType.Type.TypeName)
	}

	return nil

}

func Pokedex(config *CommandConfig, cache *pokecache.Cache, parameter string) error {

	pokedex, isInitialized := cache.Get("pokedex")

	if !isInitialized {
		fmt.Println("No pokemons caught, go ahead and catch 'em all!")
		return nil
	}

	var storedPokedex []string

	err := json.Unmarshal(pokedex, &storedPokedex)

	if err != nil {
		return err
	}

	for _, pokemon := range storedPokedex {
		fmt.Printf(" - %s\n", pokemon)
	}

	return nil

}
