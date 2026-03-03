package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/nxrmqlly/silo/internal/app"
	"github.com/nxrmqlly/silo/internal/config"
)

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		switch args[0] {
		// todo
		case "reset":
			return
		case "config":
			return
		case "help", "--help", "-h":
			return
		}
	}

	// Check if config exists
	if config.ConfigExists() {
		runSiloInteractive()
	} else {
		runSiloWizard()
	}
}

func runSiloInteractive() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	model := app.NewSiloModel(cfg.NotesDir)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func runSiloWizard() {
	fmt.Printf("Initializing SILO setup wizard...\n")

	model := app.NewWizardModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
