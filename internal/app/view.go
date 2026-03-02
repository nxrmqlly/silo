package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m CustomModel) View() tea.View {
	// s := "silo - a notes app"

	if m.editor == nil || m.statusbar == nil {
		return tea.NewView("")
	}

	editorView := m.editor.View()
	statusView := m.statusbar.View()

	return tea.NewView(lipgloss.JoinVertical(
		lipgloss.Left,
		editorView,
		statusView,
	))

}
