package app

import (
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
	"github.com/nxrmqlly/silo/internal/config"
)

var (
	// Color scheme
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true)
	stepStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	accentStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	pathStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("120")).Bold(true)
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	helpTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Italic(true)
)

type wizardStep int

const (
	stepWelcome wizardStep = iota
	stepDirInput
	stepReview
	stepSetup
	stepDone
)

type WizardModel struct {
	step        wizardStep
	textInput   textinput.Model
	configDir   string
	defaultPath string
	notesDir    string
	err         error
}

func NewWizardModel() *WizardModel {
	homeDir, _ := os.UserHomeDir()
	defaultPath := filepath.Join(homeDir, "notes")

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
				m.notesDir = path
				m.step = stepReview
			}
		}
		return m, cmd

	case stepReview:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				m.step = stepSetup
				// Perform setup
				if err := m.setup(m.notesDir); err != nil {
					m.err = err
					return m, nil
				}
				m.step = stepDone
			case "backspace":
				m.step = stepDirInput
			}
		}
		return m, nil

	case stepDone:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "enter" {
				// Transition to editor
				return NewSiloModel(m.notesDir), nil
			}
		}
		return m, nil
	}

	return m, nil
}

func (m *WizardModel) View() tea.View {
	switch m.step {

	case stepWelcome:
		return tea.NewView(
			"\n" +
				titleStyle.Render("✨ Welcome to SILO") + "\n" +
				titleStyle.Render("═══════════════════") + "\n\n" +
				"Let's set up your notes management system.\n" +
				"This wizard will help you configure SILO for the first time.\n\n" +
				helpTextStyle.Render("Press Enter to continue..."))

	case stepDirInput:
		view := "\n" +
			stepStyle.Render("Step 1 of 3") + " · " + accentStyle.Render("Choose Notes Location") + "\n" +
			subtleStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n\n" +
			"Where should your notes live?\n" +
			helpTextStyle.Render("(Press Enter to use the default)") + "\n\n"
		view += m.textInput.View()
		if m.err != nil {
			view += "\n\n" + errorStyle.Render("❌ Error: " + m.err.Error())
		}
		return tea.NewView(view)

	case stepReview:
		view := "\n" +
			stepStyle.Render("Step 2 of 3") + " · " + accentStyle.Render("Review Configuration") + "\n" +
			subtleStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n\n" +
			"Notes directory:\n" +
			"  " + pathStyle.Render(m.notesDir) + "\n\n" +
			"Config will be saved to:\n" +
			"  " + pathStyle.Render("~/.config/silo/config.json") + "\n\n" +
			helpTextStyle.Render("Backspace to edit · Enter to confirm...")
		return tea.NewView(view)

	case stepSetup:
		return tea.NewView(
			"\n" +
				stepStyle.Render("Step 3 of 3") + " · " + accentStyle.Render("Setting up...") + "\n" +
				subtleStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") + "\n\n" +
				"Creating directories and saving configuration...")

	case stepDone:
		view := "\n" +
			successStyle.Render("✓ Setup Complete") + "\n" +
			successStyle.Render("═════════════════") + "\n\n" +
			"Notes directory created at:\n" +
			"  " + pathStyle.Render(m.notesDir) + "\n\n" +
			"Config saved to:\n" +
			"  " + pathStyle.Render("~/.config/silo/config.json") + "\n\n" +
			helpTextStyle.Render("Press Enter to launch SILO...")
		return tea.NewView(view)

	}

	return tea.NewView("")
}

func (m *WizardModel) setup(notesDir string) error {
	// Create notes directory
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return err
	}

	// Save config to:
	// ~/.config/silo/config.json
	
	cfg := &config.Config{NotesDir: notesDir}
	return config.SaveConfig(cfg)
}
