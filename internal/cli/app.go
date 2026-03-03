package cli

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/nxrmqlly/silo/internal/app"
	"github.com/nxrmqlly/silo/internal/config"
	"github.com/nxrmqlly/silo/internal/wizard"
)

func Run() {
	args := os.Args[1:]

	if len(args) > 0 {
		switch args[0] {
		case "wizard":
			runSiloWizard(false)
		case "pwd":
			cmdPwd()
		case "version", "--version", "-v":
			cmdVersion()
		case "changedir":
			cmdChangeDir()
		case "reset":
			cmdReset()
		case "help", "--help", "-h":
			cmdHelp()
		default:
			fmt.Printf("unknown command: %s\n", args[0])
			fmt.Println("run `silo help` for usage")
			os.Exit(1)
		}
		return
	}
	if config.ConfigExists() {
		runSilo()
	} else {
		runSiloWizard(true)
	}
}

func runSilo() {
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

func runSiloWizard(isFirstTime bool) {
	model := wizard.NewWizardModel(isFirstTime)
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
