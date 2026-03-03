package fs

import (
	"fmt"
	"os"
	"strings"

	"github.com/nxrmqlly/silo/internal/config"
)

func InitialSetup(notesDir string) error {
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return fmt.Errorf("failed to create notes dir: %w", err)
	}

	return config.SaveConfig(&config.Config{NotesDir: notesDir})
}

func ExpandHome(path string) string {
	if !strings.HasPrefix(path, "~/") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return strings.Replace(path, "~", home, 1)
}
