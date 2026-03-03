package wizard

import (
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	// Color scheme
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true)
	stepStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("33")).Bold(true)
	accentStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	pathStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("120")).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Italic(true)
)

const siloAscii string = `
 ▄▄▄ ▄ █  ▄▄▄  
▀▄▄  ▄ █ █   █  
▄▄▄▀ █ █ ▀▄▄▄▀  
     █ █       
`
const copyNotice = `copyright (c) Ritam Das [GNU GPLv2]
https://github.com/nxrmqlly/silo`

const welcomeStr = `✨ welcome to silo - dead simple notes app for your terminal`

type WizardModel struct {
	step        wizardStep
	textInput   textinput.Model
	configDir   string
	defaultPath string
	err         error
	isFirstTime bool

	notesDir   string // set on confirm
	configPath string // set on confirm
}

type wizardStep int

const (
	stepWelcome wizardStep = iota
	stepDirInput
	stepDone
)

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

func (m *WizardModel) NotesDir() string {
	return m.notesDir
}

func (m *WizardModel) Init() tea.Cmd {
	return nil
}

func (m *WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	if key.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.step {

	case stepWelcome:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "enter" {
				m.step = stepDirInput
			}
		}

	case stepDirInput:

		switch key.String() {
		case "enter":
			path := m.textInput.Value()
			if path == "" {
				path = m.defaultPath
			}

			m.notesDir = path
			m.step = stepDone

		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

	case stepDone:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "enter" {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m *WizardModel) View() tea.View {
	var view string
	var header string

	header = lipgloss.JoinHorizontal(
		lipgloss.Center,
		titleStyle.Render(siloAscii),
		accentStyle.Render(copyNotice),
	)
	header = lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		titleStyle.Render(welcomeStr),
	)

	view += header + "\n"

	switch m.step {

	case stepWelcome:
		if m.isFirstTime {
			view += helpStyle.Render("first time run detected, let's get you set up.") + "\n"
		}
		view += helpStyle.Render("press enter to setup silo...")

	case stepDirInput:
		view += "\n" +
			stepStyle.Render("step 1/2") + " · " + accentStyle.Render("where should your notes live?") + "\n" +
			helpStyle.Render("press enter to use the default") + "\n\n" +
			m.textInput.View()

		if m.err != nil {
			view += "\n\n" + errorStyle.Render("error: "+m.err.Error())
		}

	case stepDone:
		view += "\n" +
			stepStyle.Render("step 2/2") + " · " + successStyle.Render("all set!") + "\n\n" +
			"notes dir:   " + pathStyle.Render(m.notesDir) + "\n" +
			"config file: " + pathStyle.Render(m.configPath) + "\n\n" +
			helpStyle.Render("press enter to exit · run silo again to open the editor")
	}

	return tea.NewView(view)
}
