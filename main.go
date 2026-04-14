package main

import (
    "fmt"
    "bufio"
    "os"
    "net/http"
    "log"
    "io"
    "encoding/json"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
    Next        *string
    Previous    *string
}

type PokeApi struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var commands map[string]cliCommand
var cfg *config

func init() {
    commands = map[string]cliCommand{
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    commandExit,
        },
        "help": {
            name:           "help",
            description:    "Displays a help message",
            callback:       commandHelp,
        },
        "map": {
            name:           "map",
            description:    "displays the next 20 locations in the Pokemon world",
            callback:       commandMap,
        },
        "mapb": {
            name:           "mapb",
            description:    "displays the previous 20 locations in the Pokemon world",
            callback:       commandMapb,
        },
    }

    url := "https://pokeapi.co/api/v2/location-area/"
    cfg = &config{
        Next:     &url,
        Previous: nil, 
    }

}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        text := scanner.Text()
        cleaned := cleanInput(text)
        cmd, ok := commands[cleaned[0]]
        if !ok {
            // command not found
            fmt.Println("Unknown command")
            continue
        }
        err := cmd.callback(cfg)
        if err != nil {
            fmt.Println(err)
        }
    }
}

func commandMap(cfg *config) error {

    if cfg.Next == nil {
        fmt.Println("You have traveled to far adventurer!! Try: mapb")
        return nil
    }

    res, err := http.Get(*cfg.Next)
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

    pokeapi := PokeApi{}
    jsonerr := json.Unmarshal(body, &pokeapi)
    if jsonerr != nil {
        fmt.Println(jsonerr)
    }

    fmt.Println("Next URL:", *pokeapi.Next)

    cfg.Next = pokeapi.Next
    cfg.Previous = pokeapi.Previous

    //print the names of the next 20 places from results
    for i := 0; i < len(pokeapi.Results); i++ {
        fmt.Printf("%v\n", pokeapi.Results[i].Name)
    }

    return nil
}

func commandMapb(cfg *config) error {

    if cfg.Previous == nil {
        fmt.Println("In the starting area!! Try: map")
        return nil
    }

    res, err := http.Get(*cfg.Previous)
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

    pokeapi := PokeApi{}
    jsonerr := json.Unmarshal(body, &pokeapi)
    if jsonerr != nil {
        fmt.Println(jsonerr)
    }

    cfg.Next = pokeapi.Next
    cfg.Previous = pokeapi.Previous

    //print the names of the previous 20 places from results
    for i := 0; i < len(pokeapi.Results); i++ {
        fmt.Printf("%v\n", pokeapi.Results[i].Name)
    }

    return nil
}

func commandHelp(cfg *config) error {
    fmt.Println("")
    fmt.Println("Welcome to the Pokedex!")
    fmt.Println("Usage:")
    fmt.Println("")
    for key, value := range commands {
                fmt.Printf("\t%v: %v", key, value.description)
                fmt.Println("")
            }
    fmt.Println("")
    return nil
}

func commandExit(cfg *config) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}


