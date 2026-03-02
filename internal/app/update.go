package app

import (
	"github.com/nxrmqlly/silo/internal/ui"
	tea "charm.land/bubbletea/v2"
)

func (m CustomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.editor.SetSize(m.width, m.height)
		return m, nil

	case ui.SaveFileMsg:
		// simulate save
		return m, nil
	}

	var cmd tea.Cmd
	cmd = m.editor.Update(msg)
	return m, cmd
}
