package ui

import (

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

func NewEditor() *Editor {
	ta := textarea.New()
	ta.Placeholder = "Start writing..."
	ta.ShowLineNumbers = false
	ta.Focus()

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
		Height(e.height).
		Padding(0, 1)

	if e.focused {
		style = style.Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("69"))
	} else {
		style = style.Border(lipgloss.HiddenBorder())
	}

	return style.Render(e.textarea.View())
}
