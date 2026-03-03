package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *SiloModel) View() tea.View {
	if m.editor == nil || m.statusbar == nil || m.sidebar == nil {
		return tea.NewView("")
	}

	sidebarView := m.sidebar.View()
	editorView := m.editor.View()
	statusView := m.statusbar.View()

	top := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, editorView)

	return tea.NewView(
		lipgloss.JoinVertical(lipgloss.Left, top, statusView),
	)

}
