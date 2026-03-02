package app

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/nxrmqlly/silo/internal/ui"
)

func (m *CustomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// 1:3 space
		sidebarWidth := m.width / 4
		editorWidth := m.width - sidebarWidth

		// for statusbar
		// todo: handle overflow
		contentHeight := m.height - 1

		m.editor.SetSize(editorWidth, contentHeight)
		m.sidebar.SetSize(sidebarWidth, contentHeight)
		m.statusbar.SetSize(msg.Width)

		return m, nil

	case ui.SaveFileMsg:
		// todo: false info - save logic later

		m.statusbar.SetFile(msg.Path)
		m.statusbar.SetDirty(false)

		return m, nil

	case ui.FileSelectedMsg:
		content, _ := os.ReadFile(msg.Path)
		m.editor.LoadFile(msg.Path, string(content))
	}

	cmd1 := m.editor.Update(msg)
	cmd2 := m.sidebar.Update(msg)

	// update statusbar every tick
	line, col := m.editor.CurrentCursorPosition()
	m.statusbar.SetFile(m.editor.FilePath())
	m.statusbar.SetDirty(m.editor.IsDirty())
	m.statusbar.SetCursor(line, col)
	m.statusbar.SetStats(
		m.editor.LineCount(),
		m.editor.WordCount(),
	)

	m.statusbar.SetFile(m.editor.FilePath())
	m.statusbar.SetDirty(m.editor.IsDirty())

	return m, tea.Batch(cmd1, cmd2)
}
