package ui

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Editor struct {
	textarea textarea.Model

	filePath string
	dirty    bool

	focused bool
	width   int
	height  int

	autoSave    bool
	savePending bool
}

func (e *Editor) SetAutoSave(s bool) {
	e.autoSave = s
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

func (e *Editor) CurrentContent() string {
	return e.textarea.Value()
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

func (e *Editor) LoadFile(path string, content string) {
	e.filePath = path
	e.textarea.SetValue(content)
	e.dirty = false
}

func (e *Editor) Update(msg tea.Msg) tea.Cmd {
	// Handle keyboard shortcuts
	dispatchSave := func() tea.Msg {
		return SaveFileMsg{
			Path:    e.filePath,
			Content: e.textarea.Value(),
		}
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+s":
			e.dirty = false
			e.savePending = false
			return dispatchSave
		}
	}

	// autosave timer completion
	if _, ok := msg.(AutosaveMsg); ok {
		e.savePending = false
		return dispatchSave
	}

	prev := e.textarea.Value()
	var cmd tea.Cmd

	// TODO: Add syntax highlighting for code blocks

	e.textarea, cmd = e.textarea.Update(msg)

	contentChanged := e.textarea.Value() != prev
	if contentChanged {
		e.dirty = true

		// debounce autosave
		if e.autoSave && !e.savePending {
			e.savePending = true
			return tea.Batch(
				tea.Tick(1*time.Second, func(time.Time) tea.Msg {
					return AutosaveMsg{}
				}),
				cmd,
			)
		}
	}

	return cmd
}

func (e *Editor) View() string {

	if e.width == 0 || e.height == 0 {
		return ""
	}

	style := lipgloss.NewStyle().
		Width(e.width).
		Height(e.height).
		Padding(0, 0)

	if e.focused {
		style = style.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("139"))

	} else {
		style = style.Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238"))
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
	s.Focused.CursorLineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("139"))
	s.Focused.CursorLine = lipgloss.NewStyle().Background(lipgloss.Color("236"))

	s.Blurred.CursorLineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	s.Blurred.LineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))

	ta.SetStyles(s)

	return &Editor{
		textarea: ta,
		filePath: "",
		focused:  true,
	}
}
