package app

import (
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/nxrmqlly/silo/internal/config"
)

var (
	// Color scheme
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true)
	stepStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("33")).Bold(true)
	accentStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
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

const siloAscii string = `
 ▄▄▄ ▄ █  ▄▄▄  
▀▄▄  ▄ █ █   █ 
▄▄▄▀ █ █ ▀▄▄▄▀ 
     █ █       

✨ welcome to silo - dead simple notes app for your terminal`

type WizardModel struct {
	step        wizardStep
	textInput   textinput.Model
	configDir   string
	defaultPath string
	notesDir    string
	err         error
	isFirstTime bool
}

func NewWizardModel(isFirstTime bool) *WizardModel {
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
		isFirstTime: isFirstTime,
	}
}

func (m *WizardModel) Init() tea.Cmd {
	return nil
}

func (m *WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

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
				// do the basic setup
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
				os.Exit(0)
			}
		}
		return m, nil
	}

	return m, nil
}

func (m *WizardModel) View() tea.View {
	view := titleStyle.Render(siloAscii)

	switch m.step {

	case stepWelcome:
		if m.isFirstTime {
			view += helpTextStyle.Render("\nfirst time run detected, please continue in the wizard")
		}
		view += "\n" + helpTextStyle.Render("press enter to continue...")

	case stepDirInput:
		view += "\n" +
			stepStyle.Render("step 1/3") + " · " + accentStyle.Render("where should your notes live?\n") +
			helpTextStyle.Render("(press enter to use the default)") + "\n\n"
		view += m.textInput.View()
		if m.err != nil {
			view += "\n\n" + errorStyle.Render("❌ err: "+m.err.Error())
		}

	case stepReview:
		view += "\n" +
			stepStyle.Render("step 2/3") + " · " + accentStyle.Render("review config\n") + "\n" +
			"notes directory:\n" +
			"  " + pathStyle.Render(m.notesDir) + "\n\n" +
			helpTextStyle.Render("backspace to edit · enter to confirm...")

	case stepDone:
		view += "\n" +
			stepStyle.Render("step 3/3") + " · " + successStyle.Render("setup complete") + "\n" +
			helpTextStyle.Render("press enter to exit.\nre-run silo to enter editor mode\n")
		return tea.NewView(view)

	}

	return tea.NewView(view)
}

func (m *WizardModel) setup(notesDir string) error {
	// Create notes directory
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return err
	}

	// ~/.config/silo/config.json
	cfg := &config.Config{NotesDir: notesDir}
	return config.SaveConfig(cfg)
}
