package cli

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/nxrmqlly/silo/internal/config"
	"github.com/nxrmqlly/silo/internal/fs"
)

// overriden during build time via ldflags
var Version = "dev"

const helpCmdStr = `usage: silo [command]

commands:
  (none)      open the editor
  wizard      re-run the setup wizard
  pwd         print the notes directory
  changedir   change the notes directory
  reset       delete config and start fresh
  version     print version
  help        show this message`

func cmdPwd() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("silo is not setup. run `silo` to get started.")
		os.Exit(1)
	}
	fmt.Println(cfg.NotesDir)
}

func cmdVersion() {
	fmt.Printf("silo version %s %s/%s\n", Version, runtime.GOOS, runtime.GOARCH)
}

func cmdChangeDir() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("silo is not set up. run `silo` to get started.")
		os.Exit(1)
	}

	grey := lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Italic(true)
	purp := lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))

	var prompt string = purp.Render("? new notes directory? ") +
		grey.Render("("+cfg.NotesDir+")") + "\n> "

	fmt.Print(prompt)
	input := readLine()

	if input == "" {
		fmt.Println("no change.")
		return
	}

	input = fs.ExpandHome(input)

	if err := fs.CreateDir(input); err != nil {
		fmt.Printf("error creating directory: %v\n", err)
		os.Exit(1)
	}

	cfg.NotesDir = input
	if err := config.SaveConfig(cfg); err != nil {
		fmt.Printf("error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(green.Render("notes directory changed to", input))

}

func cmdReset() {
	fmt.Print("your silo config will be deleted. your notes will not be touched. continue? (y/N)\n> ")
	if readLine() != "y" {
		fmt.Println("reset cancelled.")
		return
	}

	cfgPath, err := config.GetConfigPath()
	if err != nil {
		fmt.Printf("error resolving config path: %v\n", err)
		os.Exit(1)
	}

	if err := os.Remove(cfgPath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("error deleting config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("config removed. run `silo` to set up again.")
}

func cmdHelp() {
	fmt.Println(helpCmdStr)
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
