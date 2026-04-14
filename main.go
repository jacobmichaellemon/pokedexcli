package main

import (
    "fmt"
    "bufio"
    "os"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    stop := ""
    for stop != "exit" {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        text := scanner.Text()
        stop = text
        cleaned := cleanInput(text)
        if len(cleaned) > 0 {
            message := fmt.Sprintf("Your command was: %v", cleaned[0])
            fmt.Println(message)
        } else{
            fmt.Printf("Issue processing command: %v", text)
        }

    }
}


