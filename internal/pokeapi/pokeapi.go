package pokeapi

import (
    "fmt"
    "log"
    "io"
    "net/http"
    "encoding/json"
    "pokedexcli/internal/pokecache"
)

type locations struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type encounters struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
    } `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int `json:"base_experience"`
}

func MakeRequest(url string) []byte {
    res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil
	}
	if err != nil {
		log.Fatal(err)
    }
    return body
}

func ListLocations(cache *pokecache.Cache, url string) locations {
    fmt.Println("") //format spacing

    res, ok := cache.Get(url)
    if !ok {
        res = MakeRequest(url)
        cache.Add(url, res)
	}
    
    locations := locations{}
    jsonerr := json.Unmarshal(res, &locations)
    if jsonerr != nil {
        fmt.Println(jsonerr)
    }

    //print the names of the next 20 places from results
    for i := 0; i < len(locations.Results); i++ {
        fmt.Printf("%v\n", locations.Results[i].Name)
    }

    fmt.Println("") //format spacing

    return locations
}

func ListPokemon(cache *pokecache.Cache, url string) encounters {
    fmt.Println("") //format spacing

    res, ok := cache.Get(url)
    if !ok {
        res = MakeRequest(url)
        cache.Add(url, res)
	}
    
    encounters := encounters{}
    jsonerr := json.Unmarshal(res, &encounters)
    if jsonerr != nil {
        fmt.Println(jsonerr)
    }

    //print the names of the pokemon found on the route you explored
    for i := 0; i < len(encounters.PokemonEncounters); i++ {
        fmt.Printf("%v\n", encounters.PokemonEncounters[i].Pokemon.Name)
    }

    fmt.Println("") //format spacing

    return encounters
}

func GetPokemonBaseLevel(cache *pokecache.Cache, url string) Pokemon {
    fmt.Println("") //format spacing

    res, ok := cache.Get(url)
    if !ok {
        res = MakeRequest(url)
        cache.Add(url, res)
	}
    
    pokemon := Pokemon{}
    jsonerr := json.Unmarshal(res, &pokemon)
    if jsonerr != nil {
        fmt.Println(jsonerr)
    }

    return pokemon
}

