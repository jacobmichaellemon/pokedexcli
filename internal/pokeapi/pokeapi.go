package pokeapi

import (
    "fmt"
    "log"
    "io"
    "net/http"
    "encoding/json"
    "pokedexcli/internal/pokecache"
)

type PokeApi struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func MakeRequest(url string) []byte {
    res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
    }
    return body
}

func ListLocations(cache *pokecache.Cache, url string) PokeApi {
    fmt.Println("") //format spacing

    res, ok := cache.Get(url)
    if !ok {
        res = MakeRequest(url)
        cache.Add(url, res)
	}
    
    pokeapi := PokeApi{}
    jsonerr := json.Unmarshal(res, &pokeapi)
    if jsonerr != nil {
        fmt.Println(jsonerr)
    }

    //print the names of the next 20 places from results
    for i := 0; i < len(pokeapi.Results); i++ {
        fmt.Printf("%v\n", pokeapi.Results[i].Name)
    }

    fmt.Println("") //format spacing

    return pokeapi
}