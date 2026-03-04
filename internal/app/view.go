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
	if m.isWelcome {
		rightView = m.welcome.View()
	} else if m.isPreview {
		rightView = m.preview.View()
	} else {
		rightView = m.editor.View()
	}

	top := lipgloss.JoinHorizontal(lipgloss.Top, m.sidebar.View(), rightView)

	return tea.NewView(
		lipgloss.JoinVertical(lipgloss.Left, top, m.statusbar.View()),
	)
}