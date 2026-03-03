package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/nxrmqlly/silo/internal/app"
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



	runSiloInteractive()
}

func runSiloInteractive() {
	model := app.NewSiloModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func runSiloWizard() {
	model := app.NewWizardModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}