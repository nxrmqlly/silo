package app

import (
	tea "charm.land/bubbletea/v2"
)

func (m CustomModel) View() tea.View {
	// s := "silo - a notes app"

	if m.editor == nil {
	return tea.NewView("")
}

	return tea.NewView(m.editor.View())
}
