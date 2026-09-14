package main

import (
	"encoding/json"
	"fmt"
	"os"

	"custom-waybar/modules/countdown"
	"custom-waybar/modules/updates"
	"custom-waybar/modules/weather"
	"custom-waybar/pkg/waybar"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: custom-waybar <module> [args...]")
		fmt.Fprintln(os.Stderr, "Modules: weather, updates, countdown")
		os.Exit(1)
	}

	// Module registry — add new modules here
	registry := map[string]waybar.Module{
		"weather":   &weather.Module{},
		"updates":   &updates.Module{},
		"countdown": &countdown.Module{},
	}

	moduleName := os.Args[1]
	moduleArgs := os.Args[2:]

	mod, ok := registry[moduleName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: module %q not found\n", moduleName)
		fmt.Fprintf(os.Stderr, "Available modules: weather, updates, countdown\n")
		os.Exit(1)
	}

	result, err := mod.Run(moduleArgs)
	if err != nil {
		// Write error info as waybar output with "error" class so
		// the user can style it red in their waybar.css
		errOutput := waybar.OutputClass(
			fmt.Sprintf("⚠ %s", moduleName),
			"error",
		)
		errOutput.Tooltip = fmt.Sprintf("Error: %v", err)
		json.NewEncoder(os.Stdout).Encode(errOutput)
		os.Exit(0)
	}

	// Empty text → Waybar hides the module
	if result.Text == "" {
		os.Exit(0)
	}

	json.NewEncoder(os.Stdout).Encode(result)
}