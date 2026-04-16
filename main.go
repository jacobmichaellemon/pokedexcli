package main

import (
    "fmt"
    "bufio"
    "os"
    "time"
    "pokedexcli/internal/pokeapi"
    "pokedexcli/internal/pokecache"
    "math/rand"
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
var url string
var pokedex map[string]pokeapi.Pokemon
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
        "catch": {
            name:           "catch",
            description:    "Throw a pokeball at the pokemon passed as a parameter (i.e. catch pikachu)",
            callback:       commandCatch,
        },
        "inspect": {
            name:           "inspect",
            description:    "It takes the name of a Pokemon and prints the name, height, weight, stats and type(s) of the Pokemon.",
            callback:       commandInspect,
        },
        "pokedex": {
            name:           "pokedex",
            description:    "Prints a list of all the names of the Pokemon the user has caught.",
            callback:       commandPokedex,
        },
    }

    pokedex = map[string]pokeapi.Pokemon{}

    url = "https://pokeapi.co/api/v2/location-area/"
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
    pokeapi.ListPokemon(cache, (url + area))
    return nil
}

func commandPokedex(cfg *config, val string) error {
    fmt.Println("")
    fmt.Println("Your Pokedex:")
    for key := range pokedex {
        fmt.Printf("- %v \n", key)
    }

    fmt.Println("")
    return nil
}

func commandCatch(cfg *config, mon string) error {
    fmt.Println("")
    pokemonUrl := "https://pokeapi.co/api/v2/pokemon/" + mon
    fmt.Printf("Throwing a Pokeball at %v...", mon)
    pokemon := pokeapi.GetPokemonBaseLevel(cache, pokemonUrl)
    catchChance := rand.Intn(pokemon.BaseExperience)
    throwChance := float64(float64(catchChance)/float64(pokemon.BaseExperience))
    baseChance := float64(0.35)
    if throwChance > baseChance {
        pokedex[mon] = pokemon
        fmt.Printf("%v was caught!\n", mon)
        fmt.Println("You may now inspect it with the inspect command.")
    } else {
        fmt.Printf("%v escaped!\n", mon)
    }


    fmt.Println("")
    return nil
}

func commandInspect(cfg *config, pokemon string) error {
    fmt.Println("")
    mon, ok := pokedex[pokemon]
    if !ok {
        fmt.Println("you have not caught that pokemon")
    }
    fmt.Printf("Name: %v \n", mon.Name)
    fmt.Printf("BaseExperience: %v \n", mon.BaseExperience)
    fmt.Printf("Height: %v \n", mon.Height)
    fmt.Printf("Weight: %v \n", mon.Weight)
    fmt.Println("Stats:")
    for index, _ := range mon.Stats {
        fmt.Printf("-%v: %d\n", mon.Stats[index].Stat.Name, mon.Stats[index].BaseStat)
    }
    fmt.Println("Types:")
    for index, _ := range mon.Types {
        fmt.Printf("- %v\n", mon.Types[index].Type.Name)
    }

    fmt.Println("")

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


