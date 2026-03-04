package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/nxrmqlly/silo/internal/fs"
	"github.com/nxrmqlly/silo/internal/ui"
)

func (m *SiloModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// statusbar can handle its own clear
	m.statusbar.Update(msg)

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
				switch m.rightPane {
				case PaneEditor:
					m.editor.SetFocus(true)
				case PanePreview:
					m.preview.SetFocus(true)
				case PaneWelcome:
					m.welcome.SetFocus(true)
				}
			}

		case "ctrl+x":
			if m.rightPane == PaneEditor {
				m.rightPane = PanePreview
				m.preview.SetContent(m.editor.CurrentContent())
				m.editor.SetFocus(false)
				m.preview.SetFocus(true)
			} else if m.rightPane == PanePreview {
				m.rightPane = PaneEditor
				m.preview.SetFocus(false)
				m.editor.SetFocus(true)
			}

		case "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		sidebarWidth := m.width / 4
		editorWidth := m.width - sidebarWidth
		contentHeight := m.height - 1

		if m.rightPane == PanePreview {
			m.rightPane = PaneEditor
			m.editor.SetFocus(true)
			m.preview.SetFocus(false)
		}

		m.sidebar.SetSize(sidebarWidth, contentHeight)
		m.editor.SetSize(editorWidth, contentHeight)
		m.preview.SetSize(editorWidth, contentHeight)
		m.welcome.SetSize(editorWidth, contentHeight)
		m.statusbar.SetSize(m.width)

	case ui.SaveFileMsg:
		if err := fs.WriteFile(msg.Path, msg.Content); err != nil {
			return m, m.setStatus("err save: " + err.Error())
		}
		m.statusbar.SetFile(msg.Path)
		m.statusbar.SetDirty(false)
		return m, m.setStatus("saved")

	case ui.FileSelectedMsg:
		content, err := fs.ReadFile(msg.Path)
		if err != nil {
			return m, m.setStatus("err read: " + err.Error())
		}
		m.rightPane = PaneEditor
		m.editor.SetFocus(true)
		m.preview.SetFocus(false)
		m.welcome.SetFocus(false)
		m.editor.LoadFile(msg.Path, string(content))
		m.statusbar.SetFile(msg.Path)
		m.statusbar.SetDirty(false)
		return m, nil

	case ui.FileDeletedMsg:
		if m.editor.FilePath() == msg.Path {
			m.editor.LoadFile("", "")
			m.statusbar.SetFile("")
			m.statusbar.SetDirty(false)
			m.rightPane = PaneWelcome
			m.welcome.SetFocus(true)
			m.editor.SetFocus(false)
		}
		return m, tea.Batch(
			m.setStatus("deleted"),
			func() tea.Msg { return ui.RefreshSidebarMsg{} },
		)

	case ui.FileCreatedMsg:
		return m, tea.Batch(
			m.setStatus("created: "+msg.Path),
			func() tea.Msg { return ui.RefreshSidebarMsg{} },
		)
	}

	// ? only update component in focus.
	var cmd tea.Cmd
	switch m.focus {
	case FocusRight:
		switch m.rightPane {
		case PaneEditor:
			cmd = m.editor.Update(msg)
		case PanePreview:
			cmd = m.preview.Update(msg)
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
	return m, cmd
}
