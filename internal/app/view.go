package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *SiloModel) View() tea.View {
	if m.editor == nil || m.statusbar == nil || m.sidebar == nil {
		return tea.NewView("")
	}

	var rightView string

	switch m.rightPane {
	case PaneEditor:
		rightView = m.editor.View()
	case PaneWelcome:
		rightView = m.welcome.View()
	case PanePreview:
		rightView = m.preview.View()
	}

	top := lipgloss.JoinHorizontal(lipgloss.Top, m.sidebar.View(), rightView)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, top, m.statusbar.View()))
}
