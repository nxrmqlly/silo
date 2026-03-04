package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/nxrmqlly/silo/internal/fs"
	"github.com/nxrmqlly/silo/internal/ui"
)

func (m *SiloModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.focus == FocusRight {
				m.focus = FocusSidebar
				m.sidebar.SetFocus(true)
				m.editor.SetFocus(false)
				m.preview.SetFocus(false)
				m.welcome.SetFocus(false)
			} else {
				m.focus = FocusRight
				m.sidebar.SetFocus(false)
				m.editor.SetFocus(true)
				m.preview.SetFocus(true)
				m.welcome.SetFocus(true)
			}

		case "ctrl+x":
			m.isPreview = !m.isPreview

		case "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// 1:3 space
		sidebarWidth := m.width / 4
		editorWidth := m.width - sidebarWidth

		// for statusbar
		contentHeight := m.height - 1

		m.editor.SetSize(editorWidth, contentHeight)
		m.sidebar.SetSize(sidebarWidth, contentHeight)
		m.preview.SetSize(editorWidth, contentHeight)
		m.welcome.SetSize(editorWidth, contentHeight)
		m.statusbar.SetSize(msg.Width)

		return m, nil

	case ui.SaveFileMsg:
		if err := fs.WriteFile(msg.Path, msg.Content); err != nil {
			m.statusbar.SetStatus("err save: " + err.Error())
			return m, nil
		}

		m.statusbar.SetFile(msg.Path)
		m.statusbar.SetDirty(false)
		m.statusbar.SetStatus("saved")

		m.preview.SetContent(msg.Content)

		return m, nil

	case ui.FileSelectedMsg:
		content, err := fs.ReadFile(msg.Path)
		if err != nil {
			m.statusbar.SetStatus("err read: " + err.Error())
			return m, nil
		}

		m.isWelcome = false

		m.editor.LoadFile(msg.Path, string(content))
		m.preview.SetContent(string(content))
		m.statusbar.SetFile(msg.Path)
		m.statusbar.SetDirty(false)
		return m, nil

	case ui.FileDeletedMsg:
		if m.editor.FilePath() == msg.Path {
			m.editor.LoadFile("", "")
			m.statusbar.SetFile("")
			m.statusbar.SetDirty(false)
			m.isWelcome = true
		}
		m.statusbar.SetStatus("deleted")
		return m, func() tea.Msg {
			return ui.RefreshSidebarMsg{}
		}

	case ui.FileCreatedMsg:
		m.statusbar.SetStatus("created: " + msg.Path)

		return m, func() tea.Msg {
			return ui.RefreshSidebarMsg{}
		}
	}

	// ? only update component in focus.
	var cmd tea.Cmd
	switch m.focus {
	case FocusRight:
		if m.isPreview {
			cmd = m.preview.Update(msg)
		} else {
			cmd = m.editor.Update(msg)
		}
	case FocusSidebar:
		cmd = m.sidebar.Update(msg)
	}

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
	return m, cmd
}
