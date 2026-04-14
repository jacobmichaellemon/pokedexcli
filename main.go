package main

import (
    "fmt"
    "bufio"
    "os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

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
            return
        }
        err := cmd.callback()
        if err != nil {
            fmt.Println(err)
        }
        if len(cleaned) > 0 {
            message := fmt.Sprintf("Your command was: %v", cleaned[0])
            fmt.Println(message)
        } else{
            fmt.Printf("Issue processing command: %v", text)
        }
    }
}

func commandHelp() error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Println("Usage:")
    fmt.Println("")
    for key, value := range commands {
                fmt.Printf("%v: %v", key, value.description)
                fmt.Println("")
            }
    fmt.Println("")
    return nil
}

func commandExit() error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}


