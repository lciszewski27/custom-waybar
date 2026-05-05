package main

import (
	"encoding/json"
	"fmt"
	"os"

	"custom-waybar/modules/updates"
	"custom-waybar/modules/weather"
	"custom-waybar/pkg/waybar"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	// Rejestr modułów - tutaj dodajesz nowe funkcjonalności
	registry := map[string]waybar.Module{
		"weather": &weather.Module{},
		"updates": &updates.Module{},
	}

	moduleName := os.Args[1]
	moduleArgs := os.Args[2:]

	mod, ok := registry[moduleName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Module %s not found\n", moduleName)
		os.Exit(1)
	}

	result, err := mod.Run(moduleArgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	json.NewEncoder(os.Stdout).Encode(result)
}
