package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/nxrmqlly/silo/internal/fs"
	"github.com/nxrmqlly/silo/internal/ui"
)

func (m *SiloModel) setFocus(f FocusMode) {
	m.focus = f

	m.sidebar.SetFocus(f == FocusSidebar)
	m.editor.SetFocus(f == FocusRight && m.rightPane == PaneEditor)
	m.preview.SetFocus(f == FocusRight && m.rightPane == PanePreview)
	m.welcome.SetFocus(f == FocusRight && m.rightPane == PaneWelcome)
}

func (m *SiloModel) inSidebar(x int) bool {
	return x < m.width/4
}

func (m *SiloModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// statusbar can handle its own clear
	sbCmd := m.statusbar.Update(msg)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// inverting logic
			if m.focus == FocusRight {
				m.setFocus(FocusSidebar)
			} else {
				m.setFocus(FocusRight)
			}

		case "ctrl+x":
			switch m.rightPane {
			case PaneEditor:
				m.rightPane = PanePreview
				m.setFocus(FocusRight)
				renderCmd := m.preview.SetContent(m.editor.CurrentContent())
				spinCmd := m.statusbar.StartSpinner("rendering preview")
				return m, tea.Batch(renderCmd, spinCmd)
			case PanePreview:
				m.rightPane = PaneEditor
				m.setFocus(FocusRight)
				m.statusbar.StopSpinner()
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
		m.setFocus(FocusRight)

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
			m.setFocus(FocusRight)
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

	case ui.PreviewRenderedMsg:
		cmd := m.preview.ApplyRendered(msg)
		if !m.preview.Loading() {
			m.statusbar.StopSpinner()
		}
		return m, tea.Batch(sbCmd, cmd)

	case tea.MouseClickMsg:
		if m.inSidebar(msg.X) {
			m.setFocus(FocusSidebar)
		} else {
			m.setFocus(FocusRight)
		}

	case tea.MouseWheelMsg:
		// route scroll to the right component based on where cursor is
		if m.inSidebar(msg.X) {
			switch msg.Button {
			case tea.MouseWheelDown:
				m.sidebar.ScrollDown()
			case tea.MouseWheelUp:
				m.sidebar.ScrollUp()
			}
		}

	case ui.FileRenamedMsg:
		// if the open file was renamed, update the editor's path
		if m.editor.FilePath() == msg.OldPath {
			m.editor.LoadFile(msg.NewPath, m.editor.CurrentContent())
			m.statusbar.SetFile(msg.NewPath)
		}
		return m, tea.Batch(
			m.setStatus("renamed"),
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

	return m, tea.Batch(cmd, sbCmd)
}
