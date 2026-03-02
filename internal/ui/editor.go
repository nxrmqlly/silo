package ui

import (
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Editor struct {
	textarea textarea.Model
	filePath string
	dirty    bool
	focused  bool
	width    int
	height   int
}

func (e *Editor) CurrentCursorPosition() (int, int) {
	return e.textarea.Line(), e.textarea.Column()
}

func (e *Editor) LineCount() int {
	return len(strings.Split(e.textarea.Value(), "\n"))
}

func (e *Editor) WordCount() int {
	return len(strings.Fields(e.textarea.Value()))
}

func (e *Editor) IsDirty() bool {
	return e.dirty
}

func (e *Editor) FilePath() string {
	return e.filePath
}

func (e *Editor) SetSize(width, height int) {
	if width <= 2 || height <= 2 {
		return
	}

	e.width = width
	e.height = height
	e.textarea.SetWidth(width - 2)
	e.textarea.SetHeight(height - 2)
}

func (e *Editor) SetFocus(f bool) {
	e.focused = f
	if f {
		e.textarea.Focus()
	} else {
		e.textarea.Blur()
	}
}

func (e *Editor) Update(msg tea.Msg) tea.Cmd {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+s":
			e.dirty = false
			return func() tea.Msg {
				return SaveFileMsg{
					Path:    e.filePath,
					Content: e.textarea.Value(),
				}
			}
		}
	}

	var cmd tea.Cmd
	e.textarea, cmd = e.textarea.Update(msg)
	e.dirty = true

	return cmd
}

func (e *Editor) View() string {

	if e.width == 0 || e.height == 0 {
		return ""
	}

	style := lipgloss.NewStyle().
		Width(e.width).
		Height(e.height-1).
		Padding(0, 0)

	if e.focused {
		// style = style.
		// Border(lipgloss.RoundedBorder()).
		// BorderForeground(lipgloss.Color("244"))

	} else {
		style = style.Border(lipgloss.HiddenBorder())
	}

	return style.Render(e.textarea.View())
}

func NewEditor() *Editor {
	ta := textarea.New()
	ta.Prompt = ""
	ta.Placeholder = "Start writing..."
	ta.ShowLineNumbers = true
	ta.Focus()

	// styles
	s := textarea.DefaultStyles(true)
	s.Focused.LineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	s.Focused.CursorLineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	s.Focused.CursorLine = lipgloss.NewStyle().Background(lipgloss.Color("236"))
	ta.SetStyles(s) // default to dark styles.

	ta.SetValue(`# silo
	
This is a hardcoded note.
Start typing and press Ctrl+S to simulate saving.
`)

	return &Editor{
		textarea: ta,
		filePath: "notes/welcome.md",
		focused:  true,
	}
}
