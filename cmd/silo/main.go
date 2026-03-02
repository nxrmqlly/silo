package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/nxrmqlly/silo/internal/app"
)

func main() {
	p := tea.NewProgram(app.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
