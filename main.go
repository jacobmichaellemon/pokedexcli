package main

import (
    "fmt"
    "bufio"
    "os"
    "time"
    "pokedexcli/internal/pokeapi"
    "pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type config struct {
    Next        *string
    Previous    *string
}

var commands map[string]cliCommand
var cfg *config
var cache *pokecache.Cache
const baseTime = 20 * time.Second

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
            description:    "Displays the next 20 locations in the Pokemon world",
            callback:       commandMap,
        },
        "mapb": {
            name:           "mapb",
            description:    "Displays the previous 20 locations in the Pokemon world",
            callback:       commandMapb,
        },
        "explore": {
            name:           "explore",
            description:    "Displays all of the pokemon in an area passed as a parameter (i.e. explore viridian city)",
            callback:       commandExplore,
        },
    }

    url := "https://pokeapi.co/api/v2/location-area/"
    cfg = &config{
        Next:     &url,
        Previous: nil, 
    }

    cache = pokecache.NewCache(baseTime)
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        text := scanner.Text()
        cleaned := cleanInput(text)
        area := ""
        cmd, ok := commands[cleaned[0]]
        if len(cleaned) > 1 {
            area = cleaned[1]
        }
        if !ok {
            // command not found
            fmt.Println("Unknown command")
            continue
        }
        err := cmd.callback(cfg, area)
        if err != nil {
            fmt.Println(err)
        }
    }
}

func commandMap(cfg *config, area string) error {

    if cfg.Next == nil {
        fmt.Println("You have traveled to far adventurer!! Try: mapb")
        return nil
    }

    locations := pokeapi.ListLocations(cache, *cfg.Next)

    cfg.Next = locations.Next
    cfg.Previous = locations.Previous

    return nil
}

func commandMapb(cfg *config, area string) error {

    if cfg.Previous == nil {
        fmt.Println("In the starting area!! Try: map")
        return nil
    }

    locations := pokeapi.ListLocations(cache, *cfg.Previous)

    cfg.Next = locations.Next
    cfg.Previous = locations.Previous

    return nil
}

func commandExplore(cfg *config, area string) error {
    //pokeapi := ListLocations(baseURL + area)
    return nil
}

func commandHelp(cfg *config, area string) error {
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

func commandExit(cfg *config, area string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}


