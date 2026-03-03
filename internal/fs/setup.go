package fs

import (
	"fmt"
	"os"

	"github.com/nxrmqlly/silo/internal/config"
)

func InitialSetup(notesDir string) error {
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return fmt.Errorf("failed to create notes dir: %w", err)
	}

	return config.SaveConfig(&config.Config{NotesDir: notesDir})
}
