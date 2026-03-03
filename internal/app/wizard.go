package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type wizardStep int

const (
	stepWelcome wizardStep = iota
	stepDirInput
	stepDone
)

type WizardModel struct {
	step        wizardStep
	textInput   textinput.Model
	configDir   string
	defaultPath string
	err         error
}

type Config struct {
	RootDir string `json:"root_dir"`
}

func NewWizardModel() *WizardModel {
	cfgDir, _ := os.UserConfigDir()
	defaultPath := filepath.Join(cfgDir, "silo")

	ti := textinput.New()
	ti.Placeholder = defaultPath
	ti.SetValue(defaultPath)
	ti.Focus()

	return &WizardModel{
		step:        stepWelcome,
		textInput:   ti,
		defaultPath: defaultPath,
	}
}

func (m *WizardModel) Init() tea.Cmd {
	return nil
}

func (m *WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {

	case stepWelcome:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "enter" {
				m.step = stepDirInput
			}
		}
		return m, nil

	case stepDirInput:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)

		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "enter" {
				path := m.textInput.Value()
				if path == "" {
					path = m.defaultPath
				}

				if err := m.setup(path); err != nil {
					m.err = err
					return m, nil
				}

				m.step = stepDone
			}
		}
		return m, cmd

	case stepDone:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "enter" {
				// Transition to editor
				return NewSiloModel(), nil
			}
		}
		return m, nil
	}

	return m, nil
}

func (m *WizardModel) View() tea.View {
	switch m.step {

	case stepWelcome:
		return tea.NewView("Welcome to SILO.\n\nPress Enter to begin setup.")

	case stepDirInput:
		view := "Where should Silo live?\n\n"
		view += m.textInput.View()
		view += "\n\nPress Enter to confirm."
		if m.err != nil {
			view += fmt.Sprintf("\n\nError: %v", m.err)
		}
		return tea.NewView(view)

	case stepDone:
		return tea.NewView("Setup complete.\n\nPress Enter to enter SILO.")

	}

	return tea.NewView("")
}

func (m *WizardModel) setup(root string) error {
	notesDir := filepath.Join(root, "notes")

	// Create directories
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return err
	}

	// Write config.json
	cfg := Config{RootDir: root}
	cfgPath := filepath.Join(root, "config.json")

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cfgPath, data, 0644)
}
