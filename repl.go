package main

import "strings"

func cleanInput(text string) []string {
    cleaned := strings.ToLower(text)
    result := strings.Fields(cleaned)
    return result
}
