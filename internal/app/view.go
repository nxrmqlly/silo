package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *CustomModel) View() tea.View {
	// s := "silo - a notes app"

	if m.editor == nil || m.statusbar == nil {
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
